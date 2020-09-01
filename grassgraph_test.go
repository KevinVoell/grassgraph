package grassgraph

import (
	"io/ioutil"
	"testing"
)

func TestGetGrassGraph(t *testing.T) {
	got, err := GetGrassGraph("kevinvoell")

	if err != nil {
		t.Errorf("Received error: %s", err.Error())
	}

	ioutil.WriteFile("out.png", got, 0644)
}
