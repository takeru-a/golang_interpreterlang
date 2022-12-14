package parser

import (
	"testing"

	"github.com/takeru-a/golang_interpreterlang/ast"
	"github.com/takeru-a/golang_interpreterlang/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let footbar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil{
		t.Fatalf("ParserProgram() returned nil")
	}
	if len(program.Statements) !=3{
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		//期待する識別子
		{"x"},
		{"y"},
		{"footbar"},
	}

	for i, tt:= range tests{
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier){
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser){
	errors := p.Errors()
	if len(errors) == 0{
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors{
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool{
	// tokenがletではない場合
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	
	if !ok{
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	// 識別子の名前に関するエラー
	if letStmt.Name.Value != name{
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name{
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",name, letStmt.TokenLiteral())
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T){
	input := `
	return 5;
	return 10;
	return 993322;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3{
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",len(program.Statements))
	}

	for _, stmt := range program.Statements{
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%d",len(program.Statements))
			continue
		}
		if returnStmt.TokenLiteral() != "return"{
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}