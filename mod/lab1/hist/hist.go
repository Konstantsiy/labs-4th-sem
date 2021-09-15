package hist

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"os"
)

func DrawHistogram(values plotter.Values) error {
	p := plot.New()
	p.Title.Text = "histogram plot"
	fmt.Println("count of values: ", len(values))
	hist, err := plotter.NewHist(values, 20)
	if err != nil {
		return err
	}

	p.Add(hist)

	curDir, _ := os.Getwd()
	path := curDir+"/mod/lab1/images/"
	filename := "hist.jpg"

	err = p.Save(3*vg.Inch, 3*vg.Inch, path+filename)
	if err != nil {
		return err
	}

	return nil
}
