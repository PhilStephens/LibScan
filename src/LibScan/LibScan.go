package LibScan

import (
	"errors"
	"fmt"
	. "myUtils" //' file utils and test utils
	"os"
	"strings"
	"text/scanner"
)

const (
	//' these were for early experiments; supplanted by data/go_winapi_list.txt & data/walk_list.txt
	lf = "src/goLibrarian/goLibrarian.go"
	mf = "C:/git/WCSG_TravellingSalesmanProblem/b4c_WestCoastSalesGo/src/main.go"
	wf = "C:/git/WCSG_TravellingSalesmanProblem/b4c_WestCoastSalesGo/src/westCoastSalesGo/WestCoastSalesGo.go"
	uf = "C:/git/WCSG_TravellingSalesmanProblem/b4c_WestCoastSalesGo/src/myUtils/myUtils.go"

	//' local controls
	vb  = 5 //' vb: verbosity, where lower numbers are less verbose
	lpf = 0 //' lpf: max line to look at in each file, for early debug stages; 0 for unrestricted
	mfl = 0 //' mfl: max files to look at in each lib, for early debug stages; 0 for unrestricted
)

// func GoLibrarian gathers info on the funcs etc in listed libraries (TBD: generalize), saves info
// to an array, and (TBD) prints to a file (a work in progress as of mid-Oct 2012)
func LibScan() {

	//' for low-volume results, one file
	//ScanCode(wf, 0)
	//ScanCode(wf, lpf)

	//' for now use data from data/go_winapi_list.txt & data/walk_list.txt, loaded into an array
	//' via util FileToStrAry; each string must be preceded by "c:/git/go-winapi/" or "c:/git/walk/"
	//' TBD: generalize so libraries can be specified as GUI input or at least in a file or cmd-line

	//   /*
	GoWinApiFn_As, err := FileToStrAry("\r\n", "data/go_winapi_list.txt")
	if err != nil {
		panic(err)
	}
	WalkFn_As, err := FileToStrAry("\r\n", "data/walk_list.txt")
	if err != nil {
		panic(err)
	}
	//fmt.Println("<diag.GoLibrarian.Fn files> dummy pr ", GoWinApiFn_As[0], WalkFn_As[0])

	//' announce verbosity level (prior to file loops)
	//fmt.Println("<diag.GoLibrarian> Verbosity Level", vb)
	fmt.Println("Verbosity Level", vb)

	//' Want 2-D array of strings, but could not debug (apparently) non-deterministic behavior of append,
	//' so reluctantly sw'd to an array for each column; means 4 additional params and return values(!)
	//' ...chose LibData0 etc rather than descriptive names to emphasize their membership, ease transition
	//' to this approach, and ease converting results to 2-D form
	LibData0 := make([]string, 0)
	LibData1 := make([]string, 0)
	LibData2 := make([]string, 0)
	LibData3 := make([]string, 0)
	LibData4 := make([]string, 0)

	//' a loop for each of the filename files
	lgwa := len(GoWinApiFn_As) - 1
	for ix, fn := range GoWinApiFn_As {
		if len(fn) < 5 {
			continue
		} //' to skip blank final entry
		if (ix >= mfl) && (mfl != 0) {
			break
		} //' to limit files visited, only if mfl not zero
		//fmt.Println("<diag.GoLibrarian.GoWinApiFn> file", ix+1, "of", lgwa, "is:", fn)
		fmt.Println("  [go-winapi file", ix+1, "/", lgwa, " ", fn, "]")
		LibData0, LibData1, LibData2, LibData3, LibData4 =
			ScanCode("c:/git/go-winapi/"+fn, LibData0, LibData1, LibData2, LibData3, LibData4, lpf)
	}
	lw := len(WalkFn_As) - 1
	for ix, fn := range WalkFn_As {
		if len(fn) < 5 {
			continue
		} //' to skip blank final entry
		if (ix >= mfl) && (mfl != 0) {
			break
		} //' to limit files visited, only if mfl not zero
		//fmt.Println("<diag.GoLibrarian.WalkFn> file", ix+1, "of", lw, "is:", fn)
		fmt.Println("  [walk file", ix+1, "/", lw, " ", fn, "]")
		LibData0, LibData1, LibData2, LibData3, LibData4 =
			ScanCode("c:/git/walk/"+fn, LibData0, LibData1, LibData2, LibData3, LibData4, lpf)
	}

	//' and finally, convert the array-per-column to the 2-D array I wanted in the first place
	//' for printing and for later use (-- worth the effort: result is deterministic and non repeating)
	LibData := make([][]string, len(LibData0))
	for ix, _ := range LibData0 {
		LibData[ix] = make([]string, 5)
		LibData[ix][0] = LibData0[ix]
		LibData[ix][1] = LibData1[ix]
		LibData[ix][2] = LibData2[ix]
		LibData[ix][3] = LibData3[ix]
		LibData[ix][4] = LibData4[ix]
	}
	//' as diagnostic, using PrExpdA2s to see results
	fmt.Println("\n\n<diag.GoLibrarian> results as 2-D array\n")
	//PrExpdA2s(LibData, 9999, 99)

	//' Instead of PrExpdA2s, do inline loop; 
	for ix, _ := range LibData0 {
		fmt.Printf("%-41.41s|%4.4s  |%-8s|%-28s|%s\n",
			LibData0[ix], LibData1[ix], LibData2[ix], LibData3[ix], LibData4[ix])
	}

}

// struct iis for array of token info: line, column and string for each token
// (this may move to myUtils)
type iis struct {
	Ln  int
	Col int
	Tok string
}

//' ============================== experiment in code cross-indexing ==============================

// func ScanCode tbd; and maybe rename ScanCode, something like ScanLibFile
func ScanCode(filename string, LibDataIn0, LibDataIn1, LibDataIn2, LibDataIn3, LibDataIn4 []string, lmt int) (LibData0, LibData1, LibData2, LibData3, LibData4 []string) {
	//func ScanCode(filename string, LibDataIn [][]string, lmt int) (LibData [][]string) {
	//func ScanCode(filename string, LibData [][]string, lmt int) {
	//func ScanCode(filename string, LibData, lmt int) (codeLns int) {

	//' usage: ScanCode(filePathName, lmt) //' no longer a return value; lmt usually 0 for whole file

	//' now LibDataIn0 etc in and LibData0 etc out, so create and copy to
	lenLDI := len(LibDataIn0)
	LibData0 = make([]string, lenLDI)
	LibData1 = make([]string, lenLDI)
	LibData2 = make([]string, lenLDI)
	LibData3 = make([]string, lenLDI)
	LibData4 = make([]string, lenLDI)
	copy(LibData0, LibDataIn0)
	copy(LibData1, LibDataIn1)
	copy(LibData2, LibDataIn2)
	copy(LibData3, LibDataIn3)
	copy(LibData4, LibDataIn4)

	//' this section reads a file, obtains tokens via scanner.Scanner etc; fills array tok_Aiis

	tok_Aiis := make([]iis, 0)
	this_Aiis := make([]iis, 1)

	src, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	var s scanner.Scanner
	s.Init(src)
	tok := s.Scan()
	for tok != scanner.EOF {
		tokPos := s.Pos()
		PosLn, PosCol := 0, 0 //' for scope
		tokStr := s.TokenText()
		PosStr := tokPos.String()
		//' prep for save to array
		posAry := strings.Split(PosStr, ":")
		PosLnStr := posAry[0]
		PosColStr := posAry[1]
		_, err := fmt.Sscan(PosLnStr, &PosLn)
		if err != nil {
			fmt.Println("ScanCode.Sscan:", err)
			return
		}
		_, err = fmt.Sscan(PosColStr, &PosCol)
		if err != nil {
			fmt.Println("ScanCode.Sscan:", err)
			return
		}
		//fmt.Println("<diag.ScanCode.dmy> PosLn [", PosLn, "]; PosCol [", PosCol, "]; tokStr [", tokStr, "]")
		this_Aiis[0].Ln = PosLn
		this_Aiis[0].Col = PosCol
		this_Aiis[0].Tok = tokStr
		tok_Aiis = append(tok_Aiis, this_Aiis[0])

		if (PosLn > lmt) && (lmt > 0) {
			break
		}
		tok = s.Scan()
	}

	//' this section processes contents of tok_Aiis, using tokPttnSep & tokPttn

	ixLimit := len(tok_Aiis)
	step := 1
	ix := -1
	//' next 3 set up the rules for tokPttn; are generated by tokPttnSep from a condensed format
	relPos_Ai := []int{}       //' 0 for current position, 1 for next, -1 for previous
	relPosInv_Ai := []int{}    //' 0 for true if found, 1 for true if not found (ie 'inverted')
	relPosTok_As := []string{} //' token to compare to
	name := ""                 //' for scope
	for {
		ix += step
		if ix >= ixLimit {
			break
		}
        step = 1 //' restore default in case current token not recognized

		Ln := tok_Aiis[ix].Ln
		LnStr := fmt.Sprintf("%d", Ln)
		if (ix+1<ixLimit)&&(vb>8) {
		  fmt.Println("<diag.ScanCode.detail> step, ix, Ln, tok, next-tok:",step, ix, Ln, 
		      tok_Aiis[ix].Tok, tok_Aiis[ix+1].Tok)
		}

		//' Verbosity option: suppresses print unless vb is >= the param that follows it,
		//' but does not suppress addition to arrays

		//' pattern: named func
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func|1|1|(")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 3 {
				name = tok_Aiis[ix+1].Tok
				fmt.Println(Ln, "func", name)
				//' LibData entries: file path-name, line# in file, type of entry, entity name, other info if any.
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "func")
			LibData3 = append(LibData3, name)
			LibData4 = append(LibData4, "")

			step = 1
			continue
		}
		//' pattern: unnamed parameterless inline func
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func|1|0|(|2|0|)")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Println(Ln, "parameterless inline func")
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "func")
			LibData3 = append(LibData3, "")
			LibData4 = append(LibData4, "parameterless inline func")
			step = 3
			continue
		}
		//' pattern: unnamed parameterless inline func as a parameter of a func call; later realized may be a func 'contract'
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func|1|0|(|-3|0|,|-1|0|,")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Println(Ln, "inline func as a parameter of a func call")
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "inline func as a parameter of a func call")
			LibData3 = append(LibData3, "")
			LibData4 = append(LibData4, "(possibly contract for returned func)")
			step = 2
			continue
		}
		//' pattern: method, 1 token in parens
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func|1|0|(|3|0|)")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Printf("%d method %s of (%s)\n", Ln,
					tok_Aiis[ix+4].Tok, tok_Aiis[ix+2].Tok)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "method")
			LibData3 = append(LibData3, tok_Aiis[ix+4].Tok)
			LibData4 = append(LibData4, "of ("+tok_Aiis[ix+2].Tok+")")
			step = 4
			continue
		}
		//' pattern: method, 2 tokens in parens
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func|1|0|(|4|0|)")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			//' special case, may be (t type) or (*type)
			spacer := " "
			if strings.EqualFold(tok_Aiis[ix+2].Tok, "*") {
				spacer = ""
			}
			if vb >= 2 {
				fmt.Printf("%d method %s of (%s%s%s)\n", Ln,
					tok_Aiis[ix+5].Tok, tok_Aiis[ix+2].Tok, spacer, tok_Aiis[ix+3].Tok)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "method")
			LibData3 = append(LibData3, tok_Aiis[ix+5].Tok)
			LibData4 = append(LibData4, fmt.Sprintf("of (%s%s%s)", tok_Aiis[ix+2].Tok, spacer, tok_Aiis[ix+3].Tok))

			step = 5
			continue
		}
		//' pattern: method, 3 tokens in parens
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func|1|0|(|5|0|)")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Printf("%d method %s of (%s %s%s)\n", Ln,
					tok_Aiis[ix+6].Tok, tok_Aiis[ix+2].Tok, tok_Aiis[ix+3].Tok, tok_Aiis[ix+4].Tok)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "method")
			LibData3 = append(LibData3, tok_Aiis[ix+6].Tok)
			LibData4 = append(LibData4, fmt.Sprintf("of (%s %s%s)", tok_Aiis[ix+2].Tok, tok_Aiis[ix+3].Tok, tok_Aiis[ix+4].Tok))

			step = 6
			continue
		}
		//' func and none of previous patterns match, call it an 'oddball' until figure out new pattern
		//' pattern: any func not yet handled
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|func")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 1 {
				fmt.Println(Ln, "oddball func not currently recognized by goLibrarian/ScanCode")
				contextPr(tok_Aiis, ix, ixLimit, -3, 6)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "oddball func ")
			LibData3 = append(LibData3, "")
			LibData4 = append(LibData4, "not currently recognized by goLibrarian/ScanCode")

			step = 1
			continue
		}
		//' pattern: grouped structs, list each struct name
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|type|1|0|(|3|0|struct")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			iz, err := tokSeekMirror(tok_Aiis, ix+1, ixLimit)
			if err == nil {
				//fmt.Println(Ln, ": start of grouped type ",
				//	"declaration set of structs; struct lines & names:")
				//' w/i those lines, want to print similar to regular struct print, eg
				//' whent tok_Aiis[iy].Tok is struct:
				//' ("\t", tok_Aiis[iy].Ln, "struct ", tok_Aiis[iy-1].Tok)
				for iy := ix + 2; iy < iz; iy++ {
					if strings.EqualFold(tok_Aiis[iy].Tok, "struct") {
						if vb >= 1 {
							fmt.Println("\t", tok_Aiis[iy].Ln, "struct ", tok_Aiis[iy-1].Tok)
						} //' vb
						//' appends
						LibData0 = append(LibData0, filename)
						LibData1 = append(LibData1, LnStr)
						LibData2 = append(LibData2, "struct")
						LibData3 = append(LibData3, tok_Aiis[iy-1].Tok)
						LibData4 = append(LibData4, "(grouped)")

					} //' EqualFold
				} //' iy
			} else {
				if vb >= 1 {
					fmt.Println(Ln,
						"start of grouped type declaration set of structs (but failed to find end)")
					fmt.Println(err)
				} //' vb
				//' appends
				LibData0 = append(LibData0, filename)
				LibData1 = append(LibData1, LnStr)
				LibData2 = append(LibData2, "error, failed to find end of grouped type declaration")
				LibData3 = append(LibData3, "")
				LibData4 = append(LibData4, "(struct as part of that)")

			} //' err
			step = iz - ix
			continue
		} //' tokPttn
		//' pattern: grouped type, list as if each type defined separately
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|type|1|0|(")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			iz, err := tokSeekMirror(tok_Aiis, ix+1, ixLimit)
			//fmt.Println(Ln,
			//	"start of grouped type declaration set (exploring that content is TBD)")
			if err == nil {
				//fmt.Println("<diag.ScanCode.tokSeekMirror> found balancing token", tok_Aiis[iz].Tok, "on line", tok_Aiis[iz].Ln)
				//fmt.Println(Ln,
				//"start of grouped type declaration set (found that it spans", 
				//Ln, "to", tok_Aiis[iz].Ln, ")")
				//' rough plan: (inline, not func) tokens ix+1 to iz-1, with prepended \n whenever
				//' .Ln chgs; also, simulate what individual type stmts would have looked like
				//fmt.Println(Ln, ": start of grouped type ",
				//  "declaration set; simulating individual type stmts:")
				if vb >= 1 {
					fmt.Println(Ln, ": start grouped type ")
				} //' vb
				lnOld := Ln
				iyOld := ix + 2
				for iy := ix + 2; iy < iz; iy++ {
					//' rethinking this for array appends, created tokRange to return a string
					//' containing the tokens in a range, so can find range of previous line
					//' then generate an append for it
					lnNew := tok_Aiis[iy].Ln
					if lnNew > lnOld {
						lnOld = lnNew
						//fmt.Print("\n\ttype ")
						tokStr := tokRange(tok_Aiis, ixLimit, iyOld, iy+1)
						if vb >= 1 {
							fmt.Println("\ttype ", tokStr)
						} //' vb
						//' appends
						LibData0 = append(LibData0, filename)
						LibData1 = append(LibData1, LnStr)
						LibData2 = append(LibData2, "type")
						LibData3 = append(LibData3, tokStr)
						LibData4 = append(LibData4, "(grouped)")
						iyOld = iy + 2

					} //' lnNew > lnOld
					//fmt.Print(" ", tok_Aiis[iy].Tok)
				}
			} else {
				if vb >= 1 {
					fmt.Println(Ln,
						"start of grouped type declaration set (but failed to find end)")
					fmt.Println(err)
				} //' vb
				//' appends
				LibData0 = append(LibData0, filename)
				LibData1 = append(LibData1, LnStr)
				LibData2 = append(LibData2, "error, failed to find end of grouped type declaration")
				LibData3 = append(LibData3, "")
				LibData4 = append(LibData4, "")

			} //' err
			step = iz - ix
			continue
		}
		//' pattern: 'type xx func' is a func contract; func used as parameter of this type must satisfy this signature
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|type|2|0|func")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Println(Ln, "type", tok_Aiis[ix+1].Tok,
					"is a func contract")
				//' func used as parameter of this type must satisfy this signature
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "type")
			LibData3 = append(LibData3, tok_Aiis[ix+1].Tok)
			LibData4 = append(LibData4, "is a func contract")

			step = 3
			continue
		}
		//' pattern: solo struct declaration
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|type|2|0|struct")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Println(Ln, "struct", tok_Aiis[ix+1].Tok)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "struct")
			LibData3 = append(LibData3, tok_Aiis[ix+1].Tok)
			LibData4 = append(LibData4, "")

			step = 1
			continue
		}
		//' pattern: solo type declaration
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|type|2|1|struct")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 2 {
				fmt.Println(Ln, "type", tok_Aiis[ix+1].Tok)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "type")
			LibData3 = append(LibData3, tok_Aiis[ix+1].Tok)
			LibData4 = append(LibData4, "")

			step = 2
			continue
		}
		//' pattern: parameterless inline func deferred until close
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|defer|1|0|func|2|0|(|3|0|)")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 1 {
				fmt.Println(Ln,
					"parameterless inline func deferred until close")
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "deferred")
			LibData3 = append(LibData3, "")
			LibData4 = append(LibData4, "parameterless inline func")

			step = 4
			continue
		}
		//' pattern: inline func (not parameterless) deferred until close
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|defer|1|0|func|2|0|(|3|1|)")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			if vb >= 1 {
				fmt.Println(Ln,
					"inline func defered until close")
				//' BTW: no examples of this seen yet, included 'just in case'
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "defered func")
			LibData3 = append(LibData3, "")
			LibData4 = append(LibData4, "inline func")

			step = 4
			continue
		}
		//' pattern: solo const declaration if ( is not next
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|const|1|1|(")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			//' in order to use whole (non-cmt) rest of same line, a loop to find tok of next line
			iz := ix
			for iy := ix; iy < ixLimit; iy++ {
				if tok_Aiis[iy].Ln > Ln {
					iz = iy - 1
					break
				} //' Ln
			} //' iy
			tokStr := tokRange(tok_Aiis, ixLimit, ix+1, iz)
			if vb >= 1 {
				fmt.Println(Ln, "const", tokStr)
			}
			LibData0 = append(LibData0, filename)
			LibData1 = append(LibData1, LnStr)
			LibData2 = append(LibData2, "const")
			LibData3 = append(LibData3, tokStr)
			LibData4 = append(LibData4, "")

			step = 2
			continue
		}
		//' pattern: grouped const declaration if ( is next
		relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|const|1|0|(")
		if tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As) {
			iz, err := tokSeekMirror(tok_Aiis, ix+1, ixLimit)
			//fmt.Println(Ln,
			//  "start of grouped type declaration set (exploring that content is TBD)")
			if err == nil {
				if vb >= 1 {
					fmt.Println(Ln, "start grouped const ")
				} //' vb
				lnOld := Ln + 1
				iyOld := ix + 2
				for iy := ix + 2; iy <= iz; iy++ {
					lnNew := tok_Aiis[iy].Ln
					if lnNew > lnOld {
						lnOld = lnNew
						//fmt.Print("\n\ttype ")
						tokStr := tokRange(tok_Aiis, ixLimit, iyOld, iy-1)
						if vb >= 1 {
							fmt.Println("\tconst ", tokStr)
						} //' vb
						//' appends
						LibData0 = append(LibData0, filename)
						LibData1 = append(LibData1, LnStr)
						LibData2 = append(LibData2, "const")
						LibData3 = append(LibData3, tokStr)
						LibData4 = append(LibData4, "(grouped)")
						iyOld = iy

					} //' lnNew > lnOld
					//fmt.Print(" ", tok_Aiis[iy].Tok)
				} //' iy
			} //' err
		} //' tokPttn
		//' for now, ignore all other tokens... not flagging as unrecognized
	} //' for (loop exits only on break)
	//fmt.Println("<diag.ScanCode> just before return, example lines", LibData[0], "\n", LibData[2])

	return
} //' ScanCode

// func tokPttnSep takes first char of strIn as delimiter for remainder of string, splits remainder into an array of
// strings which it then rearranges into the 3 parallel output arrays, converting to int as needed; panics if any int
// conversion fails, or if the length of temporary array is not a multiple of 3
func tokPttnSep(strIn string) (relPos_Ai, relPosInv_Ai []int, relPosTok_As []string) {
	//' usage eg: relPos_Ai, relPosInv_Ai, relPosTok_As = tokPttnSep("|0|0|type|1|1|struct")

	//' TBD: it's silly to translate same strings over and over; low priority

	sep := strIn[0:1]
	//fmt.Println("<diag:tokPttnSep> sep is:", sep)
	tail := strIn[1:]
	tmpAry := strings.Split(tail, sep)
	tmpLen := len(tmpAry)
	ixMax := tmpLen / 3
	//fmt.Println("<diag:tokPttnSep> sep:", sep, ", tail:", tail, ", tmpLen:", tmpLen, ", ixMax:", ixMax, ", tmpAry:", tmpAry)
	if (tmpLen % 3) > 0 {
		fmt.Println("Input string of tokPttnSep must split into a multiple of 3, but <[", strIn, "]> splits into", tmpLen)
		panic("non-multiple of 3")
	}
	tmpI := 0
	relPos_Ai = make([]int, ixMax)       //' 0 for current position, 1 for next, -1 for previous
	relPosInv_Ai = make([]int, ixMax)    //' 0 for true if found, 1 for true if not found (ie 'inverted')
	relPosTok_As = make([]string, ixMax) //' token to compare to
	for ix := 0; ix < ixMax; ix++ {
		_, err := fmt.Sscan(tmpAry[3*ix], &tmpI)
		if err != nil {
			panic(err)
		}
		relPos_Ai[ix] = tmpI
		_, err = fmt.Sscan(tmpAry[3*ix+1], &tmpI)
		if err != nil {
			panic(err)
		}
		relPosInv_Ai[ix] = tmpI
		relPosTok_As[ix] = tmpAry[3*ix+2]
		//fmt.Println("<diag:tokPttnSep> ix:", ix, ", 3*ix +2:", 3*ix +2, ", tmpAry[3*ix +2]:", tmpAry[3*ix +2], ", deposited into relPosTok_As[ix]:", relPosTok_As[ix])
	} //' ix
	//fmt.Println("<diag:tokPttnSep> relPos_Ai:", relPos_Ai, ", relPosInv_Ai:", relPosInv_Ai, ", relPosTok_As:", relPosTok_As) 
	return
} //' tokPttnSep

// tokPttn regularizes the test for following and/or preceding tokens matching some 'arbitrary' pattern: tok_Aiis is an
// array of iis (int, int, string) with line number, column and token for each token; ix is the current position being
// examined; ixLimit is length of array so can avoid range errors; vb and vbt are current verbosity level and this 
// pattern's threshold; relPos_Ai, relPosInv_Ai, relPosTok_As all have same length and for each entry in each of them
// provide the relative position, whether to invert success bool, and token to look for; returns true if pattern found
func tokPttn(tok_Aiis []iis, ix, ixLimit int, relPos_Ai, relPosInv_Ai []int, relPosTok_As []string) bool {
	//' usage: tokPttn(tok_Aiis, ix, ixLimit, relPos_Ai, relPosInv_Ai, relPosTok_As)

	for iy, relPos := range relPos_Ai {
		//' for debug, peek at various values
		/*
		   fmt.Println("<diag:tokPttn> ix:", ix, ", iy:", iy, ", relPos:", relPos, 
		           ", relPos>ixLimit-ix:", relPos>ixLimit-ix, ", inv:", relPosInv_Ai[iy], 
		           ", ptok:", relPosTok_As[iy], ", tok_Aiis[ix+relPos].Tok:", tok_Aiis[ix+relPos].Tok)
		   fmt.Println("<diag:tokPttn> relPos_Ai:", relPos_Ai, ", relPosInv_Ai:", relPosInv_Ai, ", relPosTok_As:", relPosTok_As) 
		   //   */

		//' return false if any entry in relPos_Ai larger than ixLimit-ix
		//if relPos>ixLimit-ix { //' relPos = 0 for current position, 1 for next, -1 for previous
		if ix+relPos >= ixLimit { //' relPos = 0 for current position, 1 for next, -1 for previous
			return false
		}

		//' return false if any sub-match fails (used EqualFold because it's familiar)
		inv := relPosInv_Ai[iy]  //' 0 for true if found, 1 for true if not found (ie 'inverted')
		ptok := relPosTok_As[iy] //' token to compare to
		if (inv == 1) && strings.EqualFold(tok_Aiis[ix+relPos].Tok, ptok) {
			return false
		}
		if (inv == 0) && !strings.EqualFold(tok_Aiis[ix+relPos].Tok, ptok) {
			return false
		}
	}

	//' else return true
	return true
} //' func tokPttn

// func contextPr safely prints as many as possible (after checking limits) of the range of positions
// start to end (relative to ix) 
func contextPr(tok_Aiis []iis, ix, ixLimit, start, end int) {
	//' usage: contextPr(tok_Aiis, ix, ixLimit, start, end)
	ix0 := MaxInt(0, MinInt(ix+start, ixLimit-1)) //' ensure w/i range
	ix1 := MaxInt(0, MinInt(ix+end, ixLimit-1))   //' ensure w/i range
	step := SignInt(end - start)                  //' reverse stepping allowed but not recommended, might confuse

	//' Precede this func with a 'fmt.Print' to identify file if desired

	//' Note, "%+d" forces sign to print; ix0-ix & ix1-ix show whether start & end are 'corrected';
	//' and the line numbers are crucial for hinting at relevance or lack thereof
	fmt.Printf("context tokens from ix%+d (line %d) to ix%+d (line %d), stepping by %+d:\n",
		ix0-ix, tok_Aiis[ix0].Ln, ix1-ix, tok_Aiis[ix1].Ln, step)
	iy := ix0
	for {
		fmt.Print(tok_Aiis[iy].Tok, " ")
		if iy == ix1 {
			break
		}
		iy += step
	}
	fmt.Print("\n")
} //' func contextPr

// func tokRange (safely) returns as string as many as possible (after checking limits) of the range
// of positions start to end (not relative to ix) 
func tokRange(tok_Aiis []iis, ixLimit, start, end int) (str string) {
	//' usage: str := tokRange(tok_Aiis, ixLimit, start, end)
	ix0 := MaxInt(0, MinInt(start, ixLimit-1)) //' ensure w/i range
	ix1 := MaxInt(0, MinInt(end, ixLimit-1))   //' ensure w/i range
	step := SignInt(end - start)               //' reverse stepping allowed but not recommended, might confuse
	str = ""

	//' Precede this func with a 'fmt.Print' to identify file if desired

	//' Note, "%+d" forces sign to print; ix0-ix & ix1-ix show whether start & end are 'corrected';
	//' and the line numbers are crucial for hinting at relevance or lack thereof
	//fmt.Printf("context tokens from ix%+d (line %d) to ix%+d (line %d), stepping by %+d:\n",
	//    ix0-ix, tok_Aiis[ix0].Ln, ix1-ix, tok_Aiis[ix1].Ln, step)
	iy := ix0
	for {
		str += fmt.Sprint(tok_Aiis[iy].Tok, " ")
		if iy == ix1 {
			break
		}
		iy += step
	}
	//fmt.Print("\n")
	return
} //' func tokRange

// func tokSeekMirror seeks the paren (etc) to balance the one at tok_Aiis[ix].Tok [note: might not be
// same 'ix' as in loop of calling pgm]; token must be a single char, one of "()[]{}<>", and direction
// char is 'facing' determines direction of search [fwd for "([{<"].  Nested pairs of same type are
// skipped [more accurately they increment/decrement depth, and depth 0 is sought].  Returns error if
// tok_Aiis[ix].Tok does not contain an allowed char, or if runs out of tokens before finding balance.
// If successful returns position of balancing token, and err = nil.
func tokSeekMirror(tok_Aiis []iis, ix, ixLimit int) (result int, err error) {
	//' usage: result, err := tokSeekMirror(tok_Aiis, ix, ixLimit)

	//' this convoluted logic more compact than a switchcase, might not be faster though
	mirrorChars := "()[]{}<>" //' caution, < & > probably won't wk in this context
	cha := tok_Aiis[ix].Tok
	aa := strings.IndexAny(mirrorChars, cha)
	if aa < 0 {
		err = errors.New(cha + "non-mirror char")
		return
	}
	bb := (aa + 1) % 2
	cc := bb + (aa/2)*2
	chb := (string)(mirrorChars[cc])
	depth := 0
	result = 0

	//fmt.Println("dmy", chb)

	//' could implement single loop for both directions of search, but less error-prone seperately

	if bb > 0 { //' fwd search, avoid going past ixLimit -1
		for iy := ix; iy < ixLimit; iy++ {
			chc := tok_Aiis[iy].Tok //' ignore whether it's one char or more char
			if strings.EqualFold(chc, cha) {
				depth++
			}
			if strings.EqualFold(chc, chb) {
				depth--
			}
			//' when successful, break loop by doing return
			if depth == 0 {
				result = iy
				return
			}
		}
		//' if gets here failed, so report error
		err = errors.New("exhausted tokens w/o finding mirror char for" + cha)
		return
	} else { //' bkwd search, avoid going past 0
		for iy := ix; iy > 0; iy-- {
			chc := tok_Aiis[iy].Tok //' ignore whether it's one char or more char
			if strings.EqualFold(chc, cha) {
				depth++
			}
			if strings.EqualFold(chc, chb) {
				depth--
			}
			//' when successful, break loop by doing return
			if depth == 0 {
				result = iy
				return
			}
		}
		//' if gets here failed, so report error
		err = errors.New("exhausted tokens w/o finding mirror char for" + cha)
		return
	} //' bb, as direction of search
	return

} //' func tokSeekMirror
