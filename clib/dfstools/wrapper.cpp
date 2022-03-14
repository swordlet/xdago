#include <string>
#include <cstring>
#include "wrapper.h"
//#include <cstdint>
#include "dfslib_random.h"
#include "crc.h"

static int g_keylen = 0;
static struct xdnet_keys g_test_xkeys;
static struct dnet_keys *g_dnet_keys = nullptr;
static struct dnet_keys *g_dnet_user_keys = nullptr;
static struct dfslib_crypt *g_test_crypt = nullptr;
static struct dfslib_crypt *g_dnet_user_crypt = nullptr;
static struct dfslib_crypt *g_crypt = nullptr;


static void dnet_sector_to_password(uint32_t sector[SECTOR_SIZE / 4], char password[PWDLEN + 1]) {
    int i;
    for (i = 0; i < PWDLEN / 8; ++i) {
        unsigned crc = crc_of_array((unsigned char *) (sector + i * SECTOR_SIZE / 4 / (PWDLEN / 8)),
                                    SECTOR_SIZE / (PWDLEN / 8));
        sprintf(password + 8 * i, "%08X", crc);
    }
}

static void dnet_make_key(dfsrsa_t *key, int keylen) {
    unsigned i;
    for (i = keylen; i < DNET_KEYLEN; i += keylen) {
        memcpy(key + i, key, keylen * sizeof(dfsrsa_t));
    }
}

//生成random
static void dnet_random_sector(uint32_t sector[SECTOR_SIZE / 4]) {
    char password[PWDLEN + 1] = "Iyf&%d#$jhPo_t|3fgd+hf(s@;)F5D7gli^kjtrd%.kflP(7*5gt;Y1sYRC4VGL&";
    int i, j;
    for (i = 0; i < 3; ++i) {
        struct dfslib_string str;
        dfslib_utf8_string(&str, password, PWDLEN);
        dfslib_random_sector(sector, 0, &str, &str);
        for (j = KEYLEN_MIN / 8; j <= SECTOR_SIZE / 4; j += KEYLEN_MIN / 8)
            sector[j - 1] &= 0x7FFFFFFF;
        if (i == 2) break;
        dfsrsa_crypt((dfsrsa_t *) sector, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->priv.key, DNET_KEYLEN);
        dnet_sector_to_password(sector, password);
    }
}

int dnet_generate_random_array(void *array, unsigned long size) {
    uint32_t sector[SECTOR_SIZE / 4];
    unsigned long i;
    if (size < 4 || size & (size - 1)) return -1;
    if (size >= 512) {
        for (i = 0; i < size; i += 512) dnet_random_sector((uint32_t *) ((uint8_t *) array + i));
    } else {
        dnet_random_sector(sector);
        for (i = 0; i < size; i += 4) {
            *(uint32_t *) ((uint8_t *) array + i) = crc_of_array((unsigned char *) sector + i * 512 / size, 512 / size);
        }
    }
    return 0;
}

//todo add by myron
static int dnet_detect_keylen(dfsrsa_t *key, int keylen) {

    if (g_keylen && (key == g_dnet_keys->priv.key || key == g_dnet_keys->pub.key))
        return g_keylen;
    while (keylen >= 8) {
        if (memcmp(key, key + keylen / 2, keylen * sizeof(dfsrsa_t) / 2)) break;
        keylen /= 2;
    }
    return keylen;
}

static int set_user_crypt(struct dfslib_string *pwd) {
    uint32_t sector0[128];
    int i;
    g_dnet_user_crypt = (struct dfslib_crypt *) malloc(sizeof(struct dfslib_crypt));
    if (!g_dnet_user_crypt) return -1;
    //置0
    memset(g_dnet_user_crypt->pwd, 0, sizeof(g_dnet_user_crypt->pwd));
    dfslib_crypt_set_password(g_dnet_user_crypt, pwd);
    for (i = 0; i < 128; ++i) sector0[i] = 0x4ab29f51u + i * 0xc3807e6du;
    for (i = 0; i < 128; ++i) {
        dfslib_crypt_set_sector0(g_dnet_user_crypt, sector0);
        dfslib_encrypt_sector(g_dnet_user_crypt, sector0, 0x3e9c1d624a8b570full + i * 0x9d2e61fc538704abull);
    }
    return 0;
}

static int dnet_test_keys(void) {
    uint32_t src[SECTOR_SIZE / 4], dest[SECTOR_SIZE / 4];
    dnet_random_sector(src);
    memcpy(dest, src, SECTOR_SIZE);
    if (dfsrsa_crypt((dfsrsa_t *) dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->priv.key, DNET_KEYLEN)) return 1;
    if (dfsrsa_crypt((dfsrsa_t *) dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->pub.key, DNET_KEYLEN)) return 2;
    if (memcmp(dest, src, SECTOR_SIZE)) return 3;
    memcpy(dest, src, SECTOR_SIZE);
    if (dfsrsa_crypt((dfsrsa_t *) dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->pub.key, DNET_KEYLEN)) return 4;
    if (dfsrsa_crypt((dfsrsa_t *) dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->priv.key, DNET_KEYLEN)) return 5;
    if (memcmp(dest, src, SECTOR_SIZE)) return 6;
    return 0;
}

int cryptStart() {
    struct dfslib_string str;
    uint32_t sector0[128];
    int i;

    g_crypt = (struct dfslib_crypt *) malloc(sizeof(struct dfslib_crypt));
    if (!g_crypt) return -1;
    dfslib_crypt_set_password(g_crypt, dfslib_utf8_string(&str, MINERS_PWD, strlen(MINERS_PWD)));

    for (i = 0; i < 128; ++i) {
        sector0[i] = SECTOR0_BASE + i * SECTOR0_OFFSET;
    }

    for (i = 0; i < 128; ++i) {
        dfslib_crypt_set_sector0(g_crypt, sector0);
        dfslib_encrypt_sector(g_crypt, sector0, SECTOR0_BASE + i * SECTOR0_OFFSET);
    }

    return 0;
}


void dfslibEncryptArray(void *data, int nField, dfs64 sectorNo) {
    int pos = 0;
    //解密数据这里要循环解密
    for (int i = 0; i < nField; i++) {
        dfslib_encrypt_array(g_crypt, (uint32_t *) data + pos, DATA_SIZE, sectorNo++);
        pos = pos + 1;
    }
}

void dfslibDecryptArray(void *data, int nField, dfs64 sectorNo) {
    int pos = 0;
    //解密数据这里要循环解密
    for (int i = 0; i < nField; i++) {
        dfslib_uncrypt_array(g_crypt, (uint32_t *) data + pos, DATA_SIZE, sectorNo++);
        pos = pos + 1;
    }
}

int loadDnetKeys(void *keybytes, int length) {
    if (length != 3072) {
        return -1;
    }
    memcpy(&g_test_xkeys, keybytes, sizeof(struct xdnet_keys));

    return 3072;
}

int dnetCryptInit() {

    char password[PWDLEN + 1];
    struct dfslib_string str;
    g_test_crypt = (struct dfslib_crypt *) malloc(sizeof(struct dfslib_crypt));
    g_dnet_user_crypt = (struct dfslib_crypt *) malloc(sizeof(struct dfslib_crypt));
    g_dnet_keys = (struct dnet_keys *) malloc(sizeof(struct dnet_keys));

    memset(g_test_crypt, 0, sizeof(struct dfslib_crypt));
    memset(g_dnet_user_crypt, 0, sizeof(struct dfslib_crypt));
    memset(g_dnet_keys, 0, sizeof(struct dnet_keys));

    if (crc_init()) {
        return -1;
    }

    dnet_sector_to_password(g_test_xkeys.sect0.word, password);
    //为密码做置换操作
    dfslib_crypt_set_password(g_test_crypt, dfslib_utf8_string(&str, password, PWDLEN));
    //加密密码到sector0
    dfslib_crypt_set_sector0(g_test_crypt, g_test_xkeys.sect0.word);

    return 0;
}

int dfslibEncryptByteSector(void *raw, dfs64 sectorNo) {
    return dfslib_encrypt_sector(g_test_crypt, (dfs32 *) raw, sectorNo);
}

int dfslibDecryptByteSector(void *encrypted, dfs64 sectorNo) {
    return dfslib_uncrypt_sector(g_test_crypt, (dfs32 *) encrypted, sectorNo);
}

void *getUserDnetCrypt() {

    return (void *) g_dnet_user_crypt;
}

void *getDnetKeys() {

    return (void *) g_dnet_keys;
}

int setUserDnetCrypt(char *pwd) {

    struct dfslib_string str;
    uint32_t sector0[128];

    dfslib_utf8_string(&str, pwd, strlen(pwd));
    dfslib_crypt_set_password(g_dnet_user_crypt, &str);
    for (int i = 0; i < 128; ++i) {
        sector0[i] = 0x4ab29f51u + i * 0xc3807e6du;
    }

    for (int i = 0; i < 128; ++i) {
        dfslib_crypt_set_sector0(g_dnet_user_crypt, sector0);
        dfslib_encrypt_sector(g_dnet_user_crypt, sector0, 0x3e9c1d624a8b570full + i * 0x9d2e61fc538704abull);
    }

    return 0;
}

void setUserRandom(char *random_keys) {
    struct dfslib_string str;
    dfslib_random_fill(g_dnet_keys->pub.key, DNET_KEYLEN * sizeof(dfsrsa_t), 0,
                       dfslib_utf8_string(&str, random_keys, strlen(random_keys)));
}

void *makeDnetKeys(int keylen) {
    dfsrsa_keygen(g_dnet_keys->priv.key, g_dnet_keys->pub.key, keylen);
    dnet_make_key(g_dnet_keys->priv.key, keylen);
    dnet_make_key(g_dnet_keys->pub.key, keylen);

    if (g_dnet_user_crypt) {
        for (int i = 0; i < 4; ++i) {
            dfslib_encrypt_sector(g_dnet_user_crypt, (uint32_t *) g_dnet_keys + 128 * i, ~(uint64_t) i);
        }
    }


    if (g_dnet_user_crypt) {
        for (int i = 0; i < 4; ++i) {
            dfslib_uncrypt_sector(g_dnet_user_crypt, (uint32_t *) g_dnet_keys + 128 * i, ~(uint64_t) i);
        }
    }

    return (void *) g_dnet_keys;
}


int encryptWalletKey(void *privKey, dfs64 n) {
    //8 * 32
    return dfslib_encrypt_array(g_dnet_user_crypt, (uint32_t *) privKey, 8, n);
}

int decryptWalletKey(void *privKey, dfs64 n) {
    //8 * 32
    return dfslib_uncrypt_array(g_dnet_user_crypt, (uint32_t *) privKey, 8, n);
}

//generate_random_array
int generateRandomArray(void *array, uint32_t size) {
    //8*32
    return dnet_generate_random_array(array, size);
}

void dfslibRandomInit() {
    dfslib_random_init();
}

void crcInit() {
    crc_init();
}

int dnetDetectKeylen(int len) {
    return dnet_detect_keylen(g_dnet_keys->pub.key, len);
}


int verifyDnetKey(char *pwd, void *key) {

    memcpy(g_dnet_keys, key, 2048);

    struct dfslib_string str;
    uint32_t sector0[128];

    dfslib_utf8_string(&str, pwd, strlen(pwd));

    g_keylen = dnet_detect_keylen(g_dnet_keys->pub.key, DNET_KEYLEN);

    if (dnet_test_keys()) {
        //----------------------------------------------------
        memset(g_dnet_user_crypt->pwd, 0, sizeof(g_dnet_user_crypt->pwd));
        dfslib_crypt_set_password(g_dnet_user_crypt, &str);
        //给sector0进行赋值  每一个都是固定的
        for (int i = 0; i < 128; ++i) {
            sector0[i] = 0x4ab29f51u + i * 0xc3807e6du;
        }
        for (int i = 0; i < 128; ++i) {  //128次加密？？
            dfslib_crypt_set_sector0(g_dnet_user_crypt, sector0);
            dfslib_encrypt_sector(g_dnet_user_crypt, sector0, 0x3e9c1d624a8b570full + i * 0x9d2e61fc538704abull);
        }
        if (g_dnet_user_crypt) {
            for (int i = 0; i < (sizeof(struct dnet_keys) >> 9); ++i) {
                dfslib_uncrypt_sector(g_dnet_user_crypt, (uint32_t *) g_dnet_keys + 128 * i, ~(uint64_t) i);
            }
        }
        //-----------------------------------------------------------------------------------------------------
        g_keylen = 0;
        g_keylen = dnet_detect_keylen(g_dnet_keys->pub.key, DNET_KEYLEN);
    }

    return -dnet_test_keys();
}


void *generalDnetKey(char *pwd, char *random) {
    struct dfslib_string str;
    uint32_t sector0[128];


    dfslib_utf8_string(&str, pwd, strlen(pwd));
    dfslib_crypt_set_password(g_dnet_user_crypt, &str);
    for (int i = 0; i < 128; ++i) {
        sector0[i] = 0x4ab29f51u + i * 0xc3807e6du;
    }

    for (int i = 0; i < 128; ++i) {
        dfslib_crypt_set_sector0(g_dnet_user_crypt, sector0);
        dfslib_encrypt_sector(g_dnet_user_crypt, sector0, 0x3e9c1d624a8b570full + i * 0x9d2e61fc538704abull);
    }

    //rand fill
    struct dfslib_string str1;
    dfslib_random_fill(g_dnet_keys->pub.key, DNET_KEYLEN * sizeof(dfsrsa_t), 0,
                       dfslib_utf8_string(&str1, random, strlen(random)));


    dfsrsa_keygen(g_dnet_keys->priv.key, g_dnet_keys->pub.key, DNET_KEYLEN);
    dnet_make_key(g_dnet_keys->priv.key, DNET_KEYLEN);
    dnet_make_key(g_dnet_keys->pub.key, DNET_KEYLEN);

    if (g_dnet_user_crypt) {
        for (int i = 0; i < (sizeof(struct dnet_keys) >> 9); ++i) {
            dfslib_encrypt_sector(g_dnet_user_crypt, (uint32_t *) g_dnet_keys + 128 * i, ~(uint64_t) i);
        }
    }

    if (g_dnet_user_crypt) {
        for (int i = 0; i < (sizeof(struct dnet_keys) >> 9); ++i) {
            dfslib_uncrypt_sector(g_dnet_user_crypt, (uint32_t *) g_dnet_keys + 128 * i, ~(uint64_t) i);
        }
    }

    return (void *) g_dnet_keys;
}