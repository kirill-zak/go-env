package templating

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/tools/imports"

	"github.com/kirill-zak/go-env/entity"
)

// Render fills templates with provided data and writes the output to w.
func Render(
	w io.Writer, packageName string, sources []string, needImport bool, args []*entity.GetterTemplateArgs,
) error {
	headerTemplateArgs := struct {
		PackageName string
		Sources     []string
		NeedImport  bool
	}{
		PackageName: packageName,
		Sources:     sources,
		NeedImport:  needImport,
	}

	// Templates are rendered into the intermediate buffer first.
	buf := bytes.NewBufferString("")

	err := fileHeaderTemplate.Execute(buf, headerTemplateArgs)
	if err != nil {
		return fmt.Errorf("execute header template failed: %w", err)
	}

	for _, arg := range args {
		err = basicGetterTemplate.Execute(buf, arg)
		if err != nil {
			return fmt.Errorf("execute getter template for '%s' failed: %w", arg.VariableName, err)
		}
	}
	// Fix possible imports issues and reorder them.
	res, err := imports.Process("", buf.Bytes(), &imports.Options{
		TabWidth: 4,
		Comments: true,
	})
	if err != nil {
		return fmt.Errorf("process imports failed: %w", err)
	}
	// Write from intermediate to resulting buffer.
	_, err = w.Write(res)
	if err != nil {
		return fmt.Errorf("write result failed: %w", err)
	}

	return nil
}
