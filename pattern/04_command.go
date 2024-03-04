package pattern

import (
	"log"
	"os"
)

// Превращает запросы в объекты
// Плюсы: убирает прямую зависимость между объектами, вызывающими операции, и объектами, которые их выполняют;
// позволяет проще манипулировать этими запросами
// Минусы: усложняет код программы из-за введения новых классов

func ExampleCommand() {
	cf := &CommandFile{
		FileEditor: FileEditor{},
		FullPath:   "f.txt",
	}
	cs := CommandSave{
		cf,
	}
	ca := CommandAdd{
		cf,
	}
	cc := CommandClear{
		cf,
	}
	cd := CommandDelete{
		cf,
	}

	cf.Text = "Hello "
	cs.Execute()

	cf.Text = "world!"
	ca.Execute()

	cc.Execute()

	cd.Execute()
}

type Command interface {
	Execute()
}

type CommandFile struct {
	FileEditor FileEditor
	FullPath   string
	Text       string
}

type CommandSave struct {
	*CommandFile
}

func (cs *CommandSave) Execute() {
	err := cs.FileEditor.SaveToFile(cs.FullPath, cs.Text)
	if err != nil {
		log.Fatal(err)
	}
}

type CommandDelete struct {
	*CommandFile
}

func (cd *CommandDelete) Execute() {
	err := cd.FileEditor.DeleteFile(cd.FullPath)
	if err != nil {
		log.Fatal(err)
	}
}

type CommandAdd struct {
	*CommandFile
}

func (ca *CommandAdd) Execute() {
	err := ca.FileEditor.AddToFile(ca.FullPath, ca.Text)
	if err != nil {
		log.Fatal(err)
	}
}

type CommandClear struct {
	*CommandFile
}

func (cc *CommandClear) Execute() {
	err := cc.FileEditor.ClearFile(cc.FullPath)
	if err != nil {
		log.Fatal(err)
	}
}

type FileEditor struct {
}

func (fe FileEditor) SaveToFile(fullPath, text string) error {
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(text)
	return err
}

func (fe FileEditor) DeleteFile(fullPath string) error {
	err := os.Remove(fullPath)
	return err
}

func (fe FileEditor) ClearFile(fullPath string) error {
	err := os.Truncate(fullPath, 0)
	return err
}

func (fe FileEditor) AddToFile(fullPath, text string) error {
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(text)
	return err
}
