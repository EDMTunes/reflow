// Copyright 2019 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package syntax

import (
	"fmt"

	"github.com/grailbio/reflow/errors"
	"github.com/grailbio/reflow/flow"
	"github.com/grailbio/reflow/internal/scanner"
	"github.com/grailbio/reflow/types"
	"github.com/grailbio/reflow/values"
)

// CaseClause is a single case within a switch expression.
type CaseClause struct {
	// Position contains the source position of the clause.  It is set by the
	// parser.
	scanner.Position

	// Comment is the commentary text that precedes this case, if any.
	Comment string

	// Pat is the pattern of this case.  If the value of the switch matches
	// this pattern, the switch expression's value will be this cases's
	// expression's value.
	Pat *Pat

	// Expr is the expression of this case.  If the value of the switch matches
	// this cases's pattern, the switch expression's value will be this
	// expression's value.
	Expr *Expr
}

// evalSwitch is convenient context for evaluating the switch expression.
type evalSwitch struct {
	sess *Session
	env  *values.Env
	id   string

	// v is the value on which we are switching.  We will try to match this
	// against the case clause patterns.
	v values.T

	// t is the type of the value on which we are switching.
	t *types.T

	// resultT is the type of the switch expression itself (i.e. the
	// unification of the types of all the case clause expressions).
	resultT *types.T

	// pos is the position of the switch expression, given by the parser.
	pos scanner.Position
}

// idPath pairs the identifier that should be bound to the value found at the
// path with the path itself.  This just makes it a bit more convenient to pass
// Paths around with enough context to bind them in a *values.Env once they
// match.  If id == "", no value will be bound.
type idPath struct {
	id   string
	path Path
}

// String renders a tree-formatted version of c.
func (c *CaseClause) String() string {
	return fmt.Sprintf("case(%v, %v)", c.Pat, c.Expr)
}

// Equal tests whether case clause c is equivalent to case clause d.
func (c *CaseClause) Equal(d *CaseClause) bool {
	return c.Pat.Equal(d.Pat) && c.Expr.Equal(d.Expr)
}

// switchCont is the common type of the continuation functions we use.  The
// code to evaluate cases is written in a continuation-passing style, as this
// allows us to unify both immediate and deferred execution (e.g. when a Flow
// needs to be evaluated to determine which pattern matches).  When there is a
// successful match, the first argument will be true, and the Env will be the
// environment in which the case's expression should be evaluated.  Otherwise,
// the first argument will be false, and the Env must be ignored.
type switchCont func(bool, *values.Env) (values.T, error)

func (s *evalSwitch) evalCases(cs []*CaseClause) (values.T, error) {
	if len(cs) == 0 {
		// This means no case clause matched.  This is an error, as the switch
		// expression has no value.
		return nil, fmt.Errorf("%s: no case pattern matches value", s.pos)
	}
	return s.evalCase(cs[0], func(m bool, env *values.Env) (values.T, error) {
		if m {
			return s.evalExpr(cs[0], env)
		}
		// The case did not match successfully, so we keep looking.
		return s.evalCases(cs[1:])
	})
}

func (s *evalSwitch) evalCase(c *CaseClause, k switchCont) (values.T, error) {
	pattern := c.Pat
	ms := pattern.Matchers()
	ps := make([]*idPath, 0, len(ms))
	for _, m := range ms {
		p := &idPath{
			id:   m.Ident,
			path: m.Path(),
		}
		ps = append(ps, p)
	}
	env := s.env.Push()
	return s.evalPaths(ps, env, k)
}

func (s *evalSwitch) evalPaths(ps []*idPath, env *values.Env, k switchCont) (values.T, error) {
	if len(ps) == 0 {
		// If there are no more paths, then we have a successful match.
		return k(true, env)
	}
	return s.evalPath(ps[0], s.v, s.t, env,
		func(m bool, env *values.Env) (values.T, error) {
			if m {
				// The path matched, so we continue trying to match the other
				// paths.
				return s.evalPaths(ps[1:], env, k)
			}
			// The path did not match.  We are done here.
			return k(false, nil)
		})
}

func (s *evalSwitch) evalPath(
	p *idPath, v values.T, t *types.T, env *values.Env, k switchCont) (values.T, error) {
	if p.path.Done() {
		// The path matched, so we bind the value into the environment that we
		// are building up.
		if p.id != "" {
			env.Bind(p.id, v)
		}
		return k(true, env)
	}
	if f, ok := v.(*flow.Flow); ok {
		// We've hit a flow, so we can return immediately, and let the flow
		// evaluation continue our matching.
		return &flow.Flow{
			Op:         flow.K,
			Deps:       []*flow.Flow{f},
			FlowDigest: p.path.Digest(),
			K: func(vs []values.T) *flow.Flow {
				resultV, err := s.evalPath(p, vs[0], t, env, k)
				if err != nil {
					return &flow.Flow{
						Op:  flow.Val,
						Err: errors.Recover(err),
					}
				}
				return toFlow(resultV, s.resultT)
			},
		}, nil
	}
	// We throw away the error, because we know that the `Match` only uses its
	// error to describe why the match failed, which we don't care about in
	// this case.  There's some potential future world where `Match` can fail
	// in ways we want to report.  We'll need to revisit this at that time.
	nextV, nextT, ok, path, _ := p.path.Match(v, t)
	if !ok {
		return k(false, nil)
	}
	nextP := &idPath{
		id:   p.id,
		path: path,
	}
	return s.evalPath(nextP, nextV, nextT, env, k)
}

func (s *evalSwitch) evalExpr(
	c *CaseClause, env *values.Env) (values.T, error) {

	return c.Expr.eval(s.sess, env, s.id)
}

// caseUniv represents the universe of values in which the case clause patterns
// live.  To check exhaustiveness, we conceptually subtract all of the values
// matched by each pattern from this universe.  If no values remain, then the
// case patterns are exhaustive.
type caseUniv struct {
	*types.T
}

// checkCases performs static analysis on the cases of a switch expression.  We
// check two things:
//
// 1. Case exhaustiveness: are there any values of the switch expression type
// that do not match any pattern?
//
// 2. Case reachability: are there any cases that will never match any value
// (because they have already been matched by previous cases)?
//
// We do this by considering the set of possible values that cases need to
// match.  This is determined by the type of the expression being matched, t.
//
// For exhaustiveness, we check if all values are matched by some case pattern
// using the following algorithm.
//
// 1. Let V be the set of values not yet handled by a case.  This starts as the
// set of all values of t.
//
// 2. For each case, update V such that V = V - C, where C is the set of values
// that are matched by the case pattern.
//
// 3. If V is not empty, then the cases are not exhaustive, and we return an
// error.
//
// For reachability, we do something similar with the following algorithm.
//
// 1. Let V be the set of values not yet handled by a case.  This starts as the
// set of all values of t.
//
// 2. For each case, if C ∩ V = ∅, where C is the set of values that are matched
// by the case pattern, the case is unreachable.  There are no unhandled values
// that the case handles.
//
// 3. Update V such that V = V - C.
//
// We represent our sets as slices of patterns, []*Pat, e.g. [_] is the simplest
// representation of the set of all possible values of t; [[_], [_, _]]
// represents the set of list values that have one or two elements; etc.
//
// See Minus, Complement, Intersect, IntersectOneMany, and IntersectOne for the
// set-related operations used by the implementation.
func checkCases(t *types.T, pos scanner.Position, cs []*CaseClause) error {
	var caseEl errlist
	u := caseUniv{t}
	unhandled := []*Pat{{Kind: PatIgnore}}
	for _, c := range cs {
		p := c.Pat
		if len(u.IntersectOneMany(p, unhandled)) == 0 {
			// This pattern does not handle anything that is currently
			// unhandled, which means that it is redundant.
			caseEl = caseEl.Errorf(c.Position, "case is unreachable: %v", p)
		}
		unhandled = u.Minus(unhandled, p)
	}
	var el errlist
	if len(unhandled) != 0 {
		// TODO(jjc): Report example of value that is not matched by any case.
		el = el.Errorf(pos, "case patterns are not exhaustive")
	}
	// Append the case errors after the non-exhaustive error, so that the error
	// ordering is sensible.  The exhaustivity error applies to the enclosing
	// switch expression whereas the unreachability errors apply to the case
	// expressions.
	el = el.Append(caseEl.Make())
	return el.Make()
}

// Minus performs set subtraction, L - R, using the observation that
// L - R = L ∩ R∁, where R∁ is the complement of R in the universe of values U.
// For convenience of our specific implementation, rhs is given as a *Pat
// instead of []*Pat.
func (u caseUniv) Minus(lhs []*Pat, rhs *Pat) []*Pat {
	return u.Intersect(lhs, u.Complement(rhs))
}

// Intersect performs set intersection, L ∩ R, by taking union the pairwise
// intersection of the patterns in L × R.
func (u caseUniv) Intersect(lhs, rhs []*Pat) []*Pat {
	intersection := []*Pat{}
	for _, p := range lhs {
		oneIntersection := u.IntersectOneMany(p, rhs)
		intersection = append(intersection, oneIntersection...)
	}
	return intersection
}

// IntersectOneMany computes the union of the pairwise intersection of L and R,
// where L is given as a *Pat.  This is a convenience of our specific
// implementation.
func (u caseUniv) IntersectOneMany(lhs *Pat, rhs []*Pat) []*Pat {
	intersection := []*Pat{}
	for _, p := range rhs {
		oneIntersection := u.IntersectOne(lhs, p)
		if oneIntersection == nil {
			continue
		}
		intersection = append(intersection, oneIntersection)
	}
	return intersection
}

// Complement computes the complement of P in the universe U.
func (u caseUniv) Complement(p *Pat) []*Pat {
	switch p.Kind {
	case PatIdent, PatIgnore:
		return []*Pat{}
	case PatTuple:
		comp := []*Pat{}
		for i, q := range p.List {
			subU := caseUniv{u.T.Fields[i].T}
			qComp := subU.Complement(q)
			for _, r := range qComp {
				comp = append(comp, &Pat{
					Kind: PatTuple,
					List: sandwich(i, r, len(p.List)),
				})
			}
		}
		return comp
	case PatList:
		comp := make([]*Pat, len(p.List))
		// To match p, we need at least p.List elements, so the complement must
		// match shorter lists.
		for i := range p.List {
			comp[i] = &Pat{
				Kind: PatList,
				List: makeIgnoreList(i),
			}
		}
		subU := caseUniv{u.T.Elem}
		for i, q := range p.List {
			qComp := subU.Complement(q)
			for _, r := range qComp {
				comp = append(comp, &Pat{
					Kind: PatList,
					List: sandwich(i, r, len(p.List)),
				})
			}
		}
		if p.Tail == nil {
			// p has no Tail, so it does not match longer lists.  That means
			// that the complement must match longer lists.
			comp = append(comp, &Pat{
				Kind: PatList,
				List: makeIgnoreList(len(p.List) + 1),
				Tail: &Pat{Kind: PatIgnore},
			})
		} else {
			for _, c := range u.Complement(p.Tail) {
				comp = append(comp, &Pat{
					Kind: PatList,
					List: makeIgnoreList(len(p.List)),
					Tail: c,
				})
			}
		}
		return comp
	case PatStruct:
		comp := []*Pat{}
		pFields := p.FieldMap()
		for i, f := range u.T.Fields {
			fPat, ok := pFields[f.Name]
			if !ok {
				continue
			}
			subU := caseUniv{u.T.Fields[i].T}
			fComp := subU.Complement(fPat)
			for _, q := range fComp {
				comp = append(comp, &Pat{
					Kind:   PatStruct,
					Fields: []PatField{{Name: f.Name, Pat: q}},
				})
			}
		}
		return comp
	case PatVariant:
		comp := []*Pat{}
		variants := u.T.VariantMap()
		if p.Elem != nil {
			subU := caseUniv{variants[p.Tag]}
			elemComp := subU.Complement(p.Elem)
			for _, q := range elemComp {
				comp = append(comp, &Pat{
					Kind: PatVariant,
					Tag:  p.Tag,
					Elem: q,
				})
			}
		}
		for tag, elem := range variants {
			if tag == p.Tag {
				continue
			}
			var elemPat *Pat
			if elem != nil {
				elemPat = &Pat{Kind: PatIgnore}
			}
			comp = append(comp, &Pat{
				Kind: PatVariant,
				Tag:  tag,
				Elem: elemPat,
			})
		}
		return comp
	default:
		panic(fmt.Sprintf("unhandled pattern kind: %v", p.Kind))
	}
}

// IntersectOne returns the intersection of two individual *Pats.  If the
// intersection is ∅, then IntersectOne returns nil.
func (u caseUniv) IntersectOne(lhs, rhs *Pat) *Pat {
	if lhs.Kind == PatIgnore || lhs.Kind == PatIdent {
		return rhs
	}
	if rhs.Kind == PatIgnore || rhs.Kind == PatIdent {
		return lhs
	}
	if lhs.Kind != rhs.Kind {
		return nil
	}
	switch lhs.Kind {
	case PatTuple:
		list := make([]*Pat, len(lhs.List))
		for i := range lhs.List {
			subU := caseUniv{u.T.Fields[i].T}
			intersection := subU.IntersectOne(lhs.List[i], rhs.List[i])
			if intersection == nil {
				return nil
			}
			list[i] = intersection
		}
		return &Pat{
			Kind: PatTuple,
			List: list,
		}
	case PatList:
		nLhs := lhs.normalizeList()
		nRhs := rhs.normalizeList()
		// Make sure that nLhs is the shorter (or equal) length list, so
		// that the following logic can be more simple.
		if len(nRhs.list) < len(nLhs.list) {
			tmp := nLhs
			nLhs = nRhs
			nRhs = tmp
		}
		list := make([]*Pat, len(nRhs.list))
		subU := caseUniv{u.T.Elem}
		for i, rhsPat := range nRhs.list {
			var intersection *Pat
			if i < len(nLhs.list) {
				lhsPat := nLhs.list[i]
				intersection = subU.IntersectOne(lhsPat, rhsPat)
				if intersection == nil {
					// No element can match at this index, so no list can
					// match.
					return nil
				}
			} else {
				if !nLhs.allowTail {
					// nRhs is longer, and nLhs does not allow tails, so
					// the intersection is nil based on list length alone.
					return nil
				}
				// We're past the patterns in nLhs, so nLhs cannot add any
				// additional constraints on the pattern in nRhs.
				intersection = rhsPat
			}
			list[i] = intersection
		}
		var tail *Pat
		if nLhs.allowTail && nRhs.allowTail {
			tail = &Pat{Kind: PatIgnore}
		}
		return &Pat{
			Kind: PatList,
			List: list,
			Tail: tail,
		}
	case PatStruct:
		if u.T.Kind != types.StructKind {
			panic("should not have typechecked")
		}
		lhsFields := lhs.FieldMap()
		rhsFields := rhs.FieldMap()
		fields := make([]PatField, len(u.T.Fields))
		for i, f := range u.T.Fields {
			lhsPat, lhsOk := lhsFields[f.Name]
			rhsPat, rhsOk := rhsFields[f.Name]
			var intersection *Pat
			switch {
			case lhsOk && rhsOk:
				subU := caseUniv{f.T}
				intersection = subU.IntersectOne(lhsPat, rhsPat)
			case lhsOk:
				intersection = lhsPat
			case rhsOk:
				intersection = rhsPat
			default:
				intersection = &Pat{Kind: PatIgnore}
			}
			if intersection == nil {
				return nil
			}
			fields[i] = PatField{
				Name: f.Name,
				Pat:  intersection,
			}
		}
		return &Pat{
			Kind:   PatStruct,
			Fields: fields,
		}
	case PatVariant:
		if lhs.Tag != rhs.Tag {
			return nil
		}
		if lhs.Elem == nil {
			// rhs.Elem must also be nil, by virtue of pattern type-binding.
			return &Pat{Kind: PatVariant, Tag: lhs.Tag}
		}
		subU := caseUniv{u.T.VariantMap()[lhs.Tag]}
		intersection := subU.IntersectOne(lhs.Elem, rhs.Elem)
		if intersection == nil {
			return nil
		}
		return &Pat{
			Kind: PatVariant,
			Tag:  lhs.Tag,
			Elem: intersection,
		}
	default:
		return nil
	}
}

func (p *Pat) checkMatch(v values.T) bool {
	switch p.Kind {
	case PatIdent, PatIgnore:
		return true
	case PatTuple:
		tup := v.(values.Tuple)
		if len(tup) != len(p.List) {
			panic("should not have type-checked")
		}
		for i, q := range p.List {
			if !q.checkMatch(tup[i]) {
				return false
			}
		}
		return true
	case PatList:
		list := v.(values.List)
		if len(list) < len(p.List) {
			return false
		}
		if p.Tail == nil && len(p.List) < len(list) {
			return false
		}
		for i, q := range p.List {
			if !q.checkMatch(list[i]) {
				return false
			}
		}
		if p.Tail != nil {
			if !p.Tail.checkMatch(list[len(p.List):]) {
				return false
			}
		}
		return true
	case PatStruct:
		s := v.(values.Struct)
		for _, f := range p.Fields {
			if !f.Pat.checkMatch(s[f.Name]) {
				return false
			}
		}
		return true
	case PatVariant:
		variant := v.(*values.Variant)
		if variant.Tag != p.Tag {
			return false
		}
		if p.Elem == nil {
			return true
		}
		return p.Elem.checkMatch(variant.Elem)
	default:
		panic(fmt.Sprintf("unhandled pattern kind: %v", p.Kind))
	}
}

// makeIgnoreList makes a list of ignore patterns of the given length.  We use
// this in a few scenarios when generating complements (e.g. if we need to just
// match any list of a particular length).
func makeIgnoreList(length int) []*Pat {
	ps := make([]*Pat, length)
	for i := range ps {
		ps[i] = &Pat{Kind: PatIgnore}
	}
	return ps
}

func sandwich(j int, p *Pat, length int) []*Pat {
	ps := makeIgnoreList(length)
	ps[j] = p
	return ps
}

type normalizedListPat struct {
	list      []*Pat
	allowTail bool
}

// normalizeList flattens out a list pattern so that it is only composed of a
// list of patterns to match elements and a flag indicating whether it will
// accept longer lists.  This representation makes it much simpler to intersect
// list patterns.
func (p *Pat) normalizeList() normalizedListPat {
	if p.Kind != PatList {
		panic("lists only")
	}
	var (
		list    []*Pat
		currPat = p
		i       = 0
	)
	for {
		if len(currPat.List) <= i {
			tail := currPat.Tail
			if tail == nil {
				return normalizedListPat{list, false}
			}
			switch tail.Kind {
			case PatIdent, PatIgnore:
				return normalizedListPat{list, true}
			case PatList:
				currPat = tail
				i = 0
			default:
				panic("should not have typechecked")
			}
		} else {
			list = append(list, currPat.List[i])
			i++
		}
	}
}
