#ifndef X17_H
#define X17_H
#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

#include <stddef.h>
#include "sph_types.h"

int x17_test(unsigned char *pdata, const unsigned char *ptarget,
                    uint32_t nonce);
void x17_hash(void *state, const void *input);
//extern void x17_regenhash(struct work *work);
#ifdef __cplusplus
}
#endif

#endif /* X17_H */