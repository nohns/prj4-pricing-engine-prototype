package engine

type item struct {
    params ItemParams
}

func (i *item) order(qty int) {

}

func (i *item) price() float64 {
    return 0
}

func (i *item) tweakParams(param ItemParams) {

}
