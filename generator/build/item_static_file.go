package build

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/frigus02/website/generator/fs"
)

type staticFileItem struct {
	name     string
	origName string
	content  []byte

	settings *settings
}

func newStaticFileItem(file *fs.File, settings *settings) (*staticFileItem, error) {
	item := &staticFileItem{
		settings: settings,
	}
	err := item.update(file)
	return item, err
}

func (i *staticFileItem) addItem(item item) {
}

func (i *staticFileItem) isDependentOn(item item) bool {
	return false
}

func (i *staticFileItem) update(file *fs.File) error {
	var content bytes.Buffer
	hash := md5.New()
	writer := fs.NopWriteCloser(io.MultiWriter(hash, &content))

	if filepath.Ext(file.Name) == ".css" && i.settings.minifier != nil {
		writer = i.settings.minifier.Writer("text/css", writer)
	}

	_, err := io.Copy(writer, bytes.NewReader(file.Content))
	if err != nil {
		return fmt.Errorf("error copying stylesheet: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing (minify) writer for stylesheet: %v", err)
	}

	origBase := filepath.Base(file.Name)
	if origBase == ".htaccess" || origBase == "favicon.ico" {
		i.name = file.Name
	} else {
		origExt := filepath.Ext(file.Name)
		origBaseNoExt := strings.TrimSuffix(origBase, origExt)
		origPath := file.Name[0:strings.LastIndex(file.Name, "/")]

		i.name = fmt.Sprintf("%s/%s-%x%s", origPath, origBaseNoExt, hash.Sum(nil)[:8], origExt)
	}

	i.origName = file.Name
	i.content = content.Bytes()
	return nil
}

func (i *staticFileItem) render(ctx *renderContext) error {
	filename := filepath.Join(ctx.out, i.name)
	dir := filepath.Dir(filename)

	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return fmt.Errorf("error creating destination folder %s for stylesheet: %v", dir, err)
	}

	err = ioutil.WriteFile(filename, i.content, 0644)
	if err != nil {
		return fmt.Errorf("error creating new stylesheet: %v", err)
	}

	return nil
}
