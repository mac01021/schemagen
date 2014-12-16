package schemagen

import (
	"fmt"
	"strings"
)

type createStatement struct {
	header   string
	tail     string
	colDescs []string
	subTail  []string
}

func CreateStatement(tab *TableSpec) string {
	stmt := new(createStatement)
	stmt.header = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tab.Name)
	stmt.tail = "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n"
	for _, col := range tab.Columns {
		stmt.colDescs = append(stmt.colDescs, stmt.describeColumn(col)...)
	}
	stmt.describePK(tab)
	return stmt.header +
		strings.Join(stmt.colDescs, ",\n") + ",\n" +
		strings.Join(stmt.subTail, ",\n") +
		stmt.tail
}

func (stmt *createStatement) describePK(tab *TableSpec) {
	keyNames := []string{}
	for _, col := range tab.Key {
		keyNames = append(keyNames, col.Name)
	}
	stmt.subTail = append(stmt.subTail, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(keyNames, ", ")))
}

func (stmt *createStatement) describeColumn(col *ColumnSpec) (descs []string) {
	var scols []string
	var tcols []string
	for _, subcol := range subColumns(col) {
		descs = append(descs, subcol.Name+subcol.describeType()+subcol.describeSize()+subcol.describeModifiers())
		if subcol.FKTargetColumn != nil {
			scols = append(scols, subcol.Name)
			tcols = append(tcols, subcol.FKTargetColumn.Name)
		}
	}
	if col.FKTarget != nil {
		stmt.subTail = append(stmt.subTail, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
			strings.Join(scols, ", "),
			col.FKTarget.Name,
			strings.Join(tcols, ", ")))
	}
	return
}
