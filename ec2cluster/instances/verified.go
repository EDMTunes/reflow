// THIS FILE WAS AUTOMATICALLY GENERATED. DO NOT EDIT.

package instances

// VerifiedStatus captures the verification status for each Instance type.
type VerifiedStatus struct {
	// Attempted denotes whether a verification attempt has been made.
	Attempted bool
	// Verified denotes whether the instance type is verified to work for Reflow.
	Verified bool
	// ApproxETASeconds is the approximate ETA (in seconds) for Reflow to become available on this instance type.
	ApproxETASeconds int64
}

// VerifiedByRegion stores mapping of instance types to VerifiedStatus by AWS Region.
var VerifiedByRegion = make(map[string]map[string]VerifiedStatus)

func init() {
	VerifiedByRegion["us-west-2"] = map[string]VerifiedStatus{
		"c3.2xlarge":    {true, true, 73},
		"c3.4xlarge":    {true, true, 104},
		"c3.8xlarge":    {true, true, 103},
		"c3.large":      {true, true, 74},
		"c3.xlarge":     {true, true, 115},
		"c4.2xlarge":    {true, true, 62},
		"c4.4xlarge":    {true, false, 61},
		"c4.8xlarge":    {true, true, 64},
		"c4.large":      {true, true, 73},
		"c4.xlarge":     {true, true, 73},
		"c5.12xlarge":   {true, false, 62},
		"c5.18xlarge":   {true, true, 63},
		"c5.24xlarge":   {true, true, 63},
		"c5.2xlarge":    {true, true, 59},
		"c5.4xlarge":    {true, true, 61},
		"c5.9xlarge":    {true, true, 62},
		"c5.large":      {true, true, 61},
		"c5.xlarge":     {true, true, 61},
		"c5d.18xlarge":  {true, true, 59},
		"c5d.2xlarge":   {true, true, 63},
		"c5d.4xlarge":   {true, true, 70},
		"c5d.9xlarge":   {true, true, 63},
		"c5d.large":     {true, true, 62},
		"c5d.xlarge":    {true, true, 65},
		"c5n.18xlarge":  {true, false, 62},
		"c5n.2xlarge":   {true, true, 61},
		"c5n.4xlarge":   {true, false, 71},
		"c5n.9xlarge":   {true, false, 72},
		"c5n.large":     {true, false, 61},
		"c5n.xlarge":    {true, false, 71},
		"cc2.8xlarge":   {true, false, 238},
		"cr1.8xlarge":   {true, false, 0},
		"d2.2xlarge":    {true, true, 75},
		"d2.4xlarge":    {true, true, 75},
		"d2.8xlarge":    {true, true, 78},
		"d2.xlarge":     {true, true, 71},
		"f1.16xlarge":   {true, false, 72},
		"f1.2xlarge":    {true, false, 72},
		"f1.4xlarge":    {true, true, 100},
		"g2.2xlarge":    {true, true, 74},
		"g2.8xlarge":    {true, true, 78},
		"g3.16xlarge":   {true, false, 72},
		"g3.4xlarge":    {true, false, 62},
		"g3.8xlarge":    {true, false, 72},
		"g3s.xlarge":    {true, false, 72},
		"g4dn.12xlarge": {true, true, 76},
		"g4dn.16xlarge": {true, true, 68},
		"g4dn.2xlarge":  {true, true, 69},
		"g4dn.4xlarge":  {true, true, 71},
		"g4dn.8xlarge":  {true, true, 59},
		"g4dn.xlarge":   {true, true, 71},
		"h1.16xlarge":   {true, true, 71},
		"h1.2xlarge":    {true, true, 71},
		"h1.4xlarge":    {true, false, 72},
		"h1.8xlarge":    {true, true, 70},
		"hs1.8xlarge":   {true, false, 0},
		"i2.2xlarge":    {true, true, 74},
		"i2.4xlarge":    {true, true, 82},
		"i2.8xlarge":    {true, true, 77},
		"i2.xlarge":     {true, true, 70},
		"i3.16xlarge":   {true, true, 96},
		"i3.2xlarge":    {true, true, 69},
		"i3.4xlarge":    {true, true, 72},
		"i3.8xlarge":    {true, true, 77},
		"i3.large":      {true, true, 96},
		"i3.xlarge":     {true, true, 76},
		"i3en.12xlarge": {true, true, 123},
		"i3en.24xlarge": {true, true, 127},
		"i3en.2xlarge":  {true, true, 92},
		"i3en.3xlarge":  {true, true, 117},
		"i3en.6xlarge":  {true, true, 117},
		"i3en.large":    {true, true, 76},
		"i3en.xlarge":   {true, true, 88},
		"m3.2xlarge":    {true, true, 70},
		"m3.large":      {true, true, 69},
		"m3.medium":     {true, true, 107},
		"m3.xlarge":     {true, true, 73},
		"m4.10xlarge":   {true, true, 73},
		"m4.16xlarge":   {true, true, 69},
		"m4.2xlarge":    {true, true, 73},
		"m4.4xlarge":    {true, true, 58},
		"m4.large":      {true, true, 97},
		"m4.xlarge":     {true, true, 69},
		"m5.12xlarge":   {true, true, 64},
		"m5.16xlarge":   {true, false, 72},
		"m5.24xlarge":   {true, false, 71},
		"m5.2xlarge":    {true, false, 61},
		"m5.4xlarge":    {true, false, 71},
		"m5.8xlarge":    {true, false, 71},
		"m5.large":      {true, true, 64},
		"m5.xlarge":     {true, true, 61},
		"m5a.12xlarge":  {true, false, 73},
		"m5a.16xlarge":  {true, false, 72},
		"m5a.24xlarge":  {true, false, 61},
		"m5a.2xlarge":   {true, true, 61},
		"m5a.4xlarge":   {true, false, 81},
		"m5a.8xlarge":   {true, false, 60},
		"m5a.large":     {true, true, 62},
		"m5a.xlarge":    {true, true, 64},
		"m5ad.12xlarge": {true, false, 61},
		"m5ad.24xlarge": {true, false, 73},
		"m5ad.2xlarge":  {true, false, 71},
		"m5ad.4xlarge":  {true, false, 71},
		"m5ad.large":    {true, false, 61},
		"m5ad.xlarge":   {true, false, 71},
		"m5d.12xlarge":  {true, false, 72},
		"m5d.16xlarge":  {true, false, 72},
		"m5d.24xlarge":  {true, false, 71},
		"m5d.2xlarge":   {true, true, 61},
		"m5d.4xlarge":   {true, false, 71},
		"m5d.8xlarge":   {true, false, 72},
		"m5d.large":     {true, true, 65},
		"m5d.xlarge":    {true, false, 62},
		"p2.16xlarge":   {true, false, 63},
		"p2.8xlarge":    {true, true, 92},
		"p2.xlarge":     {true, true, 77},
		"p3.16xlarge":   {true, true, 101},
		"p3.2xlarge":    {true, true, 75},
		"p3.8xlarge":    {true, true, 92},
		"p3dn.24xlarge": {true, false, 297},
		"r3.2xlarge":    {true, true, 69},
		"r3.4xlarge":    {true, true, 115},
		"r3.8xlarge":    {true, true, 119},
		"r3.large":      {true, true, 71},
		"r3.xlarge":     {true, true, 58},
		"r4.16xlarge":   {true, true, 76},
		"r4.2xlarge":    {true, true, 80},
		"r4.4xlarge":    {true, true, 87},
		"r4.8xlarge":    {true, true, 77},
		"r4.large":      {true, true, 73},
		"r4.xlarge":     {true, true, 61},
		"r5.12xlarge":   {true, true, 69},
		"r5.16xlarge":   {true, true, 70},
		"r5.24xlarge":   {true, true, 76},
		"r5.2xlarge":    {true, true, 65},
		"r5.4xlarge":    {true, true, 62},
		"r5.8xlarge":    {true, true, 63},
		"r5.large":      {true, true, 65},
		"r5.xlarge":     {true, true, 63},
		"r5a.12xlarge":  {true, false, 73},
		"r5a.16xlarge":  {true, false, 60},
		"r5a.24xlarge":  {true, false, 72},
		"r5a.2xlarge":   {true, true, 66},
		"r5a.4xlarge":   {true, true, 64},
		"r5a.8xlarge":   {true, false, 72},
		"r5a.large":     {true, true, 62},
		"r5a.xlarge":    {true, true, 63},
		"r5ad.12xlarge": {true, false, 72},
		"r5ad.24xlarge": {true, false, 61},
		"r5ad.2xlarge":  {true, false, 72},
		"r5ad.4xlarge":  {true, false, 74},
		"r5ad.large":    {true, false, 62},
		"r5ad.xlarge":   {true, false, 84},
		"r5d.12xlarge":  {true, true, 62},
		"r5d.16xlarge":  {true, true, 78},
		"r5d.24xlarge":  {true, true, 67},
		"r5d.2xlarge":   {true, true, 63},
		"r5d.4xlarge":   {true, true, 63},
		"r5d.8xlarge":   {true, true, 66},
		"r5d.large":     {true, true, 62},
		"r5d.xlarge":    {true, true, 63},
		"t2.2xlarge":    {true, false, 17},
		"t2.xlarge":     {true, false, 16},
		"t3.2xlarge":    {true, true, 60},
		"t3.xlarge":     {true, true, 62},
		"t3a.2xlarge":   {true, true, 61},
		"t3a.xlarge":    {true, true, 64},
		"x1.16xlarge":   {true, false, 72},
		"x1.32xlarge":   {true, false, 61},
		"x1e.16xlarge":  {true, false, 18},
		"x1e.2xlarge":   {true, false, 17},
		"x1e.32xlarge":  {true, false, 24},
		"x1e.4xlarge":   {true, false, 17},
		"x1e.8xlarge":   {true, false, 17},
		"x1e.xlarge":    {true, false, 16},
		"z1d.12xlarge":  {true, true, 69},
		"z1d.2xlarge":   {true, false, 61},
		"z1d.3xlarge":   {true, true, 61},
		"z1d.6xlarge":   {true, true, 65},
		"z1d.large":     {true, true, 63},
		"z1d.xlarge":    {true, true, 61},
	}

}
