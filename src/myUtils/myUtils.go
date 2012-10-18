// contains copy-paste from earlier projects; most only as needed, except misc math utilities
package myUtils 

import (
    "bytes"
    "fmt"
    "io"
    "math"
    "os"
    "strings"
)

//' ============================== misc math utilities ==============================

// func MaxInt is trivial, just to avoid multiple casts required by use of math.Max on pair of int,
// [eg: c := int(math.Max(float64(a), float64(b))) -- but maybe casts are more efficient?]
func MaxInt(a, b int) (c int) {
    c = a
    if c < b {
        c = b
    }
    return
}

// func MinInt is trivial, just to avoid multiple casts required by use of math.Min on pair of int,
// [eg: c := int(math.Min(float64(a), float64(b))) -- but maybe casts are more efficient?]
func MinInt(a, b int) (c int) {
    c = a
    if c > b {
        c = b
    }
    return
}

// ResLmt limits resolution to 1/res; res should be a power of 10, eg 1000; no error checking
func ResLmt(f_in, res float64) (f_out float64) {
    return math.Trunc(f_in*res+0.499) / res //' round rather than simply truncate
} //' ResLmt

// SignInt gives sign of int as int 1 or -1; 0 is assumed positive
func SignInt(i_in int) (i_out int) {
    i_out = 1
    if i_in < 0 {
        i_out = -1
    }
    return
} //' SignInt

// SignFloat gives sign of float64 as float64 1.0 or -1.0; 0 is arbitrarily assumed positive
func SignFloat(f_in float64) (f_out float64) {
    f_out = 1
    if f_in < 0 {
        f_out = -1
    }
    return
} //' SignFloat

// AbsInt gives absolute value of an int w/o converting to and from float
func AbsInt(i_in int) (i_out int) {
    i_out = i_in
    if i_in < 0 {
        i_out = -i_in
    }
    return
} //' AbsInt

//' ================== file utilities ==================

// func FileToStr is a utility to open a file, read it, return whole file as one big String
func FileToStr(filename string) (result string, err error) {
    //' usage: str, err := FileToStr(filename)

    f, err := os.Open(filename)
    if err != nil {
        return //' will return empty string, and err
    }
    defer f.Close()

    var buf bytes.Buffer
    _, err = io.Copy(&buf, f)
    if err != nil {
        return //' will return empty string, and err
    }
    //' only if successful, copy to result
    result = buf.String()
    return
} //' FileToStr

// FileToStrAry calls FileToStr utility to read file, then uses ".split" to load it into 1D array
func FileToStrAry(sep, fileNameStr string) (result []string, err error) {
    //' usage: strAry, err := FileToStrAry(sep, fileNameStr)

    tmp, err := FileToStr(fileNameStr)
    if err != nil {
        return //' will return empty string, and err
    }
    result = strings.Split(tmp, sep)
    return
} //' FileToStrAry

//' ================== verbosity control utility ==================

// vbp takes verbosity level and target level as params, and if former is >= latter runs the func
// param (usually a Print, Println or Printf) -- abandoned, insufficient flexibility
//func vbp(vb, vbt int, pfunc func(...)) {
//    //' usage: vbp(vb, 2, fmt.Println("<diag> something here"))

//    if vb >= vbt {
//        pfunc()
//    }
//} //' vbp

//' ============================== unit test print-expecteds utilities ==============================
//' Some utils mainly for printing values to use as expected arrays (or for seeing intermediate results)

// PrExpdA2s: Convenience util to print contents of array of [][]string, pre-fmtd for copy-paste as unit-test expected
func PrExpdA2s(A2s [][]string, lim1, lim2 int) {
    //' usage: PrExpdA2s(A2s, 9999, 99)

    fmt.Println("{ ")
    for ix1, Af := range A2s {
        if (ix1 > 0) && (ix1%lim1 == 0) {
            fmt.Print("\n")
        }
        fmt.Print("\t{ ")
        for ix2, ff := range Af {
            if (ix2 > 0) && (ix2%lim2 == 0) {
                fmt.Print("\n\t")
            }
            fmt.Print(ff, ", ")
        }
        fmt.Println("},")
    }
    fmt.Println("}")
} //' PrExpdA2s
