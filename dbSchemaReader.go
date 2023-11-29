package dbSchemaReader

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Table_Struct struct {
	Table_name          string
	Table_Columns       []TableColumns
	IndexDetails        []Index_Name_Details
	ForeignKeys			[]ForeignKey_Details
	OutputFileName      string
	FunctionSignature   string
	FunctionSignature2  string
  }
  
  type Index_Name_Details struct {
	IndexName       string
	IndexColumn     []string
  }
  
  type ForeignKey_Details struct {
	FK_Column       			string
	FK_Related_Table			string
	FK_Related_Table_Column		string
  }

  
  type TableColumns struct {
	Column_name     string
	PrimaryFlag     bool
	UniqueFlag      bool
	ColumnType      string
	ColumnNameParams string
  }
  
func ReadSchema(filePath string)  []Table_Struct {

	var tableX []Table_Struct
	var table Table_Struct
	var tabColumns TableColumns
	readFile, err := os.Open(filePath)
	if err != nil {
	  fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
	  res1 := strings.Split(fileScanner.Text(), " ")
	  if len(res1) > 1 {
		if res1[0] == "CREATE" && res1[1] == "TABLE" {
		  table.Table_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
		  table.FunctionSignature2 = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
		  if strings.TrimSpace(table.Table_name[len(table.Table_name)-3:]) == `ies` {
			table.OutputFileName = strings.TrimSpace(table.Table_name[:len(table.Table_name)-3])+"y"
		  }else if strings.TrimSpace(table.Table_name[len(table.Table_name)-1:]) == `s` {
			table.OutputFileName = strings.TrimSpace(table.Table_name[:len(table.Table_name)-1])
		  }else {
			table.OutputFileName = table.Table_name
		  }
		  if strings.TrimSpace(table.Table_name[len(table.Table_name)-3:]) == `ies` {
			table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-3]+"y")
			} else if strings.TrimSpace(table.Table_name[len(table.Table_name)-1:]) == `s` {
			  table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:len(table.Table_name)-1])
			} else {
			  table.FunctionSignature = strings.ToUpper(strings.TrimSpace(table.Table_name[0:1]))+strings.TrimSpace(table.Table_name[1:])
		  }
		  table.Table_Columns = nil
		}
		if res1[0] == "" && res1[1] == "" && strings.TrimSpace(res1[2][0:1]) == `"` {
		  tabColumns.Column_name = strings.TrimSpace(res1[2][1:len(res1[2])-1])
		  tabColumns.ColumnType = strings.TrimSpace(res1[3][0:])
		  ///////////////////
		  column_name_slice := strings.Split(tabColumns.Column_name,"_")
		  // fmt.Println(column_name_slice)
		  for k := 0; k < len(column_name_slice); k++ {
			if column_name_slice[k] == "id" {
			  column_name_slice[k] = strings.ToUpper(strings.TrimSpace(column_name_slice[k]))
			}else{
			  column_name_slice[k] = strings.ToUpper(strings.TrimSpace(column_name_slice[k][0:1]))+strings.TrimSpace(column_name_slice[k][1:])
	  
			}
		  }
		  tabColumns.ColumnNameParams = strings.Join(column_name_slice,"")
		  //////////////////
		  if len(res1) > 4 {
			if res1[4] == `PRIMARY` {
			  tabColumns.PrimaryFlag = true
			} else{
			  tabColumns.PrimaryFlag = false
			}
			if res1[4] == `UNIQUE` {
			  tabColumns.UniqueFlag = true
			} else{
			  tabColumns.UniqueFlag = false            
			}
		  } else {
			tabColumns.PrimaryFlag = false
			tabColumns.UniqueFlag = false
		  }
		  table.Table_Columns = append(table.Table_Columns, tabColumns)
		}
		if res1[0] == "CREATE" && res1[1] == "INDEX" {
		  for i:=0; i<len(tableX); i++{
			if tableX[i].Table_name == strings.TrimSpace(res1[3][1:len(res1[3])-1]) { 
			  var index Index_Name_Details
			  index.IndexName = strings.TrimSpace(res1[3][1:len(res1[3])-1]) + strconv.Itoa(rand.Intn(90000))
			  for m :=4; m<len(res1); m++ {            
				indexColumnName := res1[m]
				if strings.TrimSpace(indexColumnName[0:1]) == `(` {
				  indexColumnName = strings.TrimSpace(indexColumnName[2:len(indexColumnName)-1])
				} else if strings.TrimSpace(indexColumnName[0:1]) == `"` {
				  indexColumnName = strings.TrimSpace(indexColumnName[1:len(indexColumnName)-1])
				}
				if strings.TrimSpace(indexColumnName[len(indexColumnName)-1:]) == `)` {
				  indexColumnName = strings.TrimSpace(indexColumnName[0:len(indexColumnName)-2])
				} else if strings.TrimSpace(indexColumnName[len(indexColumnName)-1:]) == `"` {
				  indexColumnName = strings.TrimSpace(indexColumnName[0:len(indexColumnName)-1])
				}
				// fmt.Println(indexColumnName)
				index.IndexColumn = append(index.IndexColumn,   indexColumnName)
			  }
			  tableX[i].IndexDetails = append(tableX[i].IndexDetails, index)
			}
		  }
		}
		if res1[0] == "ALTER" && res1[4] == "FOREIGN" {
			for i:=0; i<len(tableX); i++{
				fmt.Println("I have found ALTER FOREIGN key line")
				var fkDetails ForeignKey_Details
				if tableX[i].Table_name == strings.TrimSpace(res1[2][1:len(res1[3])-1]) { 
				  //var index Index_Name_Details
				  //index.IndexName = strings.TrimSpace(res1[3][1:len(res1[3])-1]) + strconv.Itoa(rand.Intn(90000))
				  fkDetails.FK_Column = res1[5]
				  fkDetails.FK_Column = strings.TrimSpace(fkDetails.FK_Column[2:len(fkDetails.FK_Column)-2])
				  fkDetails.FK_Related_Table = res1[7]
				  fkDetails.FK_Related_Table = strings.TrimSpace(fkDetails.FK_Related_Table[1:len(fkDetails.FK_Related_Table)-2])
				  fkDetails.FK_Related_Table_Column = res1[8]
				  fkDetails.FK_Related_Table_Column = strings.TrimSpace(fkDetails.FK_Related_Table_Column[2:len(fkDetails.FK_Related_Table_Column)-4])
				}
				tableX[i].ForeignKeys = append(tableX[i].ForeignKeys, fkDetails)
				fmt.Println(fkDetails.FK_Column, fkDetails.FK_Related_Table, fkDetails.FK_Related_Table_Column)
			}			
			
		}
	  }
	  if len(res1) == 1 {
		if res1[0] == ");" {
		  tableX = append(tableX, table)
		}
	  }
	}
	for i:=0; i<len(tableX); i++{
	//   fmt.Println("table Name: ", tableX[i].Table_name, "OutputFileName: ", tableX[i].OutputFileName,  "FunctionSignature: ", tableX[i].FunctionSignature, "FunctionSignature2: ", tableX[i].FunctionSignature2)
	  for j:=0; j<len(tableX[i].Table_Columns); j++{
		// fmt.Println("    column name: ", tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].Table_Columns[j].PrimaryFlag, tableX[i].Table_Columns[j].UniqueFlag, tableX[i].Table_Columns[j].ColumnNameParams)
	  }
	  for j:=0; j<len(tableX[i].IndexDetails); j++{
		// fmt.Println("    index name: ", tableX[i].IndexDetails[j].IndexName)
		for k:=0; k<len(tableX[i].IndexDetails[j].IndexColumn); k++{
		//   fmt.Println("    index column name: ", tableX[i].IndexDetails[j].IndexColumn[k])
		}
	  }
	  for j:=0; j<len(tableX[i].ForeignKeys); j++{
		// fmt.Println("    FK_Column: ", tableX[i].ForeignKeys[j].FK_Column, "FK_Related_Table: ", tableX[i].ForeignKeys[j].FK_Related_Table, "FK_Related_Table_Column: ", tableX[i].ForeignKeys[j].FK_Related_Table_Column)
	  }
	}
  return tableX
}