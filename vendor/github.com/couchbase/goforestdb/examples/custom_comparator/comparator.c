#include <libforestdb/forestdb.h>
#include "comparator.h"

extern int CompareBytesReversed(void *key1, size_t keylen1, void *key2, size_t keylen2);

int cmp_variable(void *key1, size_t keylen1, void *key2, size_t keylen2) {
	return CompareBytesReversed(key1, keylen1, key2, keylen2);
}
