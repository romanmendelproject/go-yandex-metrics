package crypto

import (
	"fmt"
)

func Example() {
	var metrics []byte
	key := "test"

	fmt.Println(GetHash(metrics, key))
	// Output:
	// rXEUjHnyGrnuxR6lx90rZoeS98DTU0rmayL3HGFSP7M=
}
