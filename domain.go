package schemagen

import (
	"strconv"
)

type SchemaSpec map[string]TableSpec

type TableSpec struct {
	Name    string
	Key     []*ColumnSpec
	Columns map[string]*ColumnSpec
}

type ColType int

const (
	INTEGER ColType = iota
	STRING
	FK
	TIMESTAMP
	DATE
	BINARY
	BOOLEAN
)

type ColumnSpec struct {
	Name           string
	IsPK           bool
	IsNullable     bool
	Type           ColType
	Size           int
	FKTarget       *TableSpec
	fkTargetName   string
	FKTargetColumn *ColumnSpec
}

func (col *ColumnSpec) describeType() string {
	if col.Type == BOOLEAN {
		return " BOOL"
	}
	if col.Type == INTEGER {
		return " INT"
	}
	if col.Type == STRING {
		return " VARCHAR"
	}
	if col.Type == BINARY {
		return " BINARY"
	}
	if col.Type == DATE {
		return " DATE"
	}
	if col.Type == TIMESTAMP {
		return " TIMESTAMP"
	}
	panic("cannot describe invalid type")
}

func (col *ColumnSpec) describeSize() string {
	if col.Size < 1 {
		return ""
	}
	return "(" + strconv.Itoa(col.Size) + ")"
}

func (col *ColumnSpec) describeModifiers() string {
	if col.IsNullable {
		return " NULL"
	}
	return "  NOT NULL"
}

func subColumns(col *ColumnSpec) (subcols []*ColumnSpec) {
	if col.Type != FK {
		subcols = []*ColumnSpec{col}
		return
	}
	targ := col.FKTarget
	for _, tCol := range targ.Key {
		tsubcols := subColumns(tCol)
		for _, tsubcol := range tsubcols {
			sname := col.Name + "_" + tsubcol.Name
			subcols = append(subcols, &ColumnSpec{sname, col.IsPK, col.IsNullable,
				tsubcol.Type, tsubcol.Size, targ, targ.Name,
				tsubcol})
		}
	}
	return
}
