#include "dfslib_crypt.h"
#include "dfsrsa.h"

#define SECTOR0_BASE           0x1947f3acu
#define SECTOR0_OFFSET         0x82e9d1b5u
#define BLOCK_HEADER_WORD      0x3fca9e2bu
#define MINERS_PWD             "minersgonnamine"
#define DATA_SIZE              8 //(sizeof(struct xdag_field) / sizeof(uint32_t))
#define WORKERNAME_HEADER_WORD 0xf46b9853u

#define DNET_KEY_SIZE    4096
#define DNET_KEYLEN    ((DNET_KEY_SIZE * 2) / (sizeof(dfsrsa_t) * 8))

#define SECTOR_LOG  9
#define SECTOR_SIZE (1 << SECTOR_LOG)

struct dnet_key {
    dfsrsa_t key[DNET_KEYLEN];
};

#define PWDLEN        64
#define KEYLEN_MIN    (DNET_KEYLEN / 4)

struct dnet_packet_header {
    uint8_t type;
    uint8_t ttl;
    uint16_t length;
    uint32_t crc32;
};

struct xsector {
    union {
        uint8_t byte[SECTOR_SIZE];
        uint32_t word[SECTOR_SIZE / sizeof(uint32_t)];
        struct xsector *next;
        struct dnet_packet_header head;
    };
};

struct xdnet_keys {
    struct dnet_key priv, pub;
    struct xsector sect0_encoded, sect0;
};

struct dnet_keys {
    struct dnet_key priv;
    struct dnet_key pub;
};

#ifdef __cplusplus
extern "C" {
#endif

extern int cryptStart();
extern void dfslibEncryptArray(void *data, int nField, dfs64 sectorNo);
extern void dfslibDecryptArray(void *data, int nField, dfs64 sectorNo);
extern int dnetCryptInit();
extern int loadDnetKeys(void *keybytes, int length);
extern int dfslibEncryptByteSector(void *raw, dfs64 sectorNo);
extern int dfslibDecryptByteSector(void *encrypted, dfs64 sectorNo);
extern int encryptWalletKey(void *privKey, dfs64 n);
extern int decryptWalletKey(void *privKey, dfs64 n);
extern void dfslibRandomInit();
extern void crcInit();
extern int verifyDnetKey(char *pwd, void *key);
extern void *generalDnetKey(char *pwd, char *random);

#ifdef __cplusplus
};
#endif