package codegen

import "fmt"

// BBLEntry is Single entry in BBL List
type BBLEntry struct {
	Block []IRIns
	Info  []map[*SymbolTableRegisterEntry]UseInfo
}

func (b *BBLEntry) String() string {
	repr := "\n<BBL Begin>\n"
	for _, ins := range b.Block {
		repr += fmt.Sprintf("%s\n", ins.String())
	}
	repr += "<BBL End>\n\n"
	return repr
}

// UseInfo stores life and next Use Information of a variable
type UseInfo struct {
	Live    bool
	NextUse int
}

// BBLList is the list of all the Basic Blocks
var BBLList []BBLEntry

var symbolInfo = make(map[*SymbolTableRegisterEntry]UseInfo)

func setSymbolInfo(arg SymbolTableEntry) {
	if arg != nil {
		if arg, isRegister := arg.(*SymbolTableRegisterEntry); isRegister {
			symbolInfo[arg] = UseInfo{true, -1}
		}
	}
}

// GenBBLList takes the IRCode (list of IRIns) as input & creates list of basic blocks
func GenBBLList(IRCode []IRIns) {
	if len(IRCode) == 0 {
		return
	}
	prevIndex := 0
	for index, ins := range IRCode {
		setSymbolInfo(ins.Arg1)
		setSymbolInfo(ins.Arg2)
		setSymbolInfo(ins.Dst)
		if ins.Typ == LBL && index != prevIndex {
			bbl := BBLEntry{Block: IRCode[prevIndex:index]}
			bbl = addUseInfo(bbl)
			BBLList = append(BBLList, bbl)
			prevIndex = index
		} else if ins.Typ == CBR || ins.Typ == JMP || ins.Typ == KEY || index == len(IRCode)-1 {
			bbl := BBLEntry{Block: IRCode[prevIndex : index+1]}
			bbl = addUseInfo(bbl)
			BBLList = append(BBLList, bbl)
			if index != len(IRCode)-1 {
				prevIndex = index + 1
			}
		}
	}
}

func isRegister(entry SymbolTableEntry) (*SymbolTableRegisterEntry, bool) {
	if entry != nil {
		entry, ok := entry.(*SymbolTableRegisterEntry)
		return entry, ok
	}
	return nil, false
}

// Adds Operands' UseInfo in the BBL
func addUseInfo(bbl BBLEntry) BBLEntry {
	bbl.Info = make([]map[*SymbolTableRegisterEntry]UseInfo, len(bbl.Block))
	infomap := make(map[*SymbolTableRegisterEntry]UseInfo)
	for i := len(bbl.Block) - 1; i >= 0; i-- {
		if dst, isReg := isRegister(bbl.Block[i].Dst); isReg {
			infomap[dst] = symbolInfo[dst]
			symbolInfo[dst] = UseInfo{false, -1}
		}
		if arg1, isReg := isRegister(bbl.Block[i].Arg1); isReg {
			infomap[arg1] = symbolInfo[arg1]
			symbolInfo[arg1] = UseInfo{true, i}
		}
		if arg2, isReg := isRegister(bbl.Block[i].Arg2); isReg {
			infomap[arg2] = symbolInfo[arg2]
			symbolInfo[arg2] = UseInfo{true, i}
		}
		ninfomap := make(map[*SymbolTableRegisterEntry]UseInfo)
		for k, v := range infomap {
			ninfomap[k] = v
		}
		bbl.Info[i] = ninfomap
	}
	return bbl
}