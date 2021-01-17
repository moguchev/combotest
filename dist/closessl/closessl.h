#include <stddef.h>
#include <stdint.h>

/*
*  Encrypts provided bytes buffer using a state of the art encryption
*  algorithm. Provided bytes are encrypted in-place, so buf contains result
*  of the encryption.
*
*  N.B. Performs CPU-extensive operations.
*
*  @param buf buffer with bytes to be encrypted
*  @param len number of uint8_t elements inside a buffer
*/
void closessl_encrypt(uint8_t *buf, size_t len);

/*
*  Decrypts provided bytes buffer encrypted with closessl_encrypt. Provided
*  bytes are decrypted in-place, so buf contains the result.
*
*  ClosedSSL decryption guaranteed to work fast on any input, regardless of
*  its encryption speed.
*
*  @param buf buffer with bytes to be encrypted
*  @param len number of uint8_t elements inside a buffer
*/
void closessl_decrypt(uint8_t *buf, size_t len);
