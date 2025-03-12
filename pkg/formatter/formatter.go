package formatter

type Formatter struct {
	folder           string
	useTokenStart    bool
	skipNextLineFunc bool
}

func New(folder string) *Formatter {
	return &Formatter{folder: folder}
}

func (f *Formatter) Format() error {
	files, err := getFiles(f.folder, ".zx")
	if err != nil {
		return err
	}
	for _, file := range files {
		fileFmt := NewFileFmt(file)
		err := fileFmt.Format()
		if err != nil {
			return err
		}
	}
	return nil
}
