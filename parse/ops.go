package parse

type OpType int

const (
	XFX OpType = iota
	FX
	XFY
	FY
	YFX
	XF
	YF
)

type Op map[OpType]int

// defaultOps is a set of ops totally stolen from SWI, I have
// literally no idea what 90% of these do.
var defaultOps = map[string]Op{
	"-->":                   {XFX: 1200},
	":-":                    {XFX: 1200, FX: 1200},
	"?-":                    {FX: 1200},
	"dynamic":               {FX: 1150},
	"discontiguous":         {FX: 1150},
	"initialization":        {FX: 1150},
	"meta_predicate":        {FX: 1150},
	"module_transparent":    {FX: 1150},
	"multifile":             {FX: 1150},
	"public":                {FX: 1150},
	"thread_local":          {FX: 1150},
	"thread_initialization": {FX: 1150},
	"volatile":              {FX: 1150},
	";":                     {XFY: 1100},
	"|":                     {XFY: 1100},
	"->":                    {XFY: 1050},
	"*->":                   {XFY: 1050},
	"*":                     {XFY: 1000, YFX: 400},
	":=":                    {XFX: 990},
	"\\+":                   {FY: 900},
	"<":                     {XFX: 700},
	"=":                     {XFX: 700},
	"=..":                   {XFX: 700},
	"=@=":                   {XFX: 700},
	"\\=@=":                 {XFX: 700},
	"=:=":                   {XFX: 700},
	"=<":                    {XFX: 700},
	"==":                    {XFX: 700},
	"=\\=":                  {XFX: 700},
	">":                     {XFX: 700},
	">=":                    {XFX: 700},
	"@<":                    {XFX: 700},
	"@=<":                   {XFX: 700},
	"@>":                    {XFX: 700},
	"@>=":                   {XFX: 700},
	"\\=":                   {XFX: 700},
	"\\==":                  {XFX: 700},
	"as":                    {XFX: 700},
	"is":                    {XFX: 700},
	">:<":                   {XFX: 700},
	":<":                    {XFX: 700},
	":":                     {XFY: 600},
	"+":                     {YFX: 500, FY: 200},
	"-":                     {YFX: 500, FY: 200},
	"/\\":                   {YFX: 500},
	"\\/":                   {YFX: 500},
	"xor":                   {YFX: 500},
	"?":                     {FX: 500},
	"/":                     {YFX: 400},
	"//":                    {YFX: 400},
	"div":                   {YFX: 400},
	"rdiv":                  {YFX: 400},
	"<<":                    {YFX: 400},
	">>":                    {YFX: 400},
	"mod":                   {YFX: 400},
	"rem":                   {YFX: 400},
	"**":                    {XFX: 200},
	"^":                     {XFY: 200},
	"\\":                    {FY: 200},
	".":                     {YFX: 100},
	"$":                     {FX: 1},
}
