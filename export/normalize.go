package export

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
)

// CleanupResourceFields ...
func CleanupResourceFields(in map[string][]byte) (map[string][]byte, error) {
	cleanedUpRsrcs := make(map[string][]byte)
	toDelete := []string{
		"metadata.annotations.control-plane.alpha.kubernetes.io/leader",
		"metadata.annotations.kubectl.kubernetes.io/last-applied-configuration",
		"metadata.creationTimestamp",
		"metadata.generation",
		"metadata.resourceVersion",
		"metadata.selfLink",
		"metadata.uid",
		"spec.clusterIP",
		"status",
	}

	for name, rsrc := range in {
		jsonParsed, _ := gabs.ParseJSON(rsrc)

		child, isArray := jsonParsed.S("items").Data().([]interface{})
		if !isArray {
			return nil, fmt.Errorf("Expected an items array")
		}
		newItems := make([]interface{}, 0, len(child))
		for _, child := range child {
			var ele *gabs.Container
			ele, _ = gabs.Consume(child)
			if name == "secret" {
				// do not export service account tokens
				if ele.S("type").Data().(string) != "kubernetes.io/service-account-token" {
					newItems = append(newItems, child)
				}
			} else {
				newItems = append(newItems, child)
			}
			jsonParsed.Set(newItems, "items")

			// delete metadata from items array
			for _, path := range toDelete {
				splitPath := strings.SplitN(path, ".", 3)
				if ele.Exists(splitPath...) {
					if err := ele.Delete(splitPath...); err != nil {
						return nil, err
					}
				}
			}
		}
		// after cleanup, there are no objects for this
		// resource so no need to export in result
		if jsonParsed.S("items").String() != "[]" {
			cleanedUpRsrcs[name] = jsonParsed.Bytes()
		}
	}

	return cleanedUpRsrcs, nil
}
