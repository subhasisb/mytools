#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "tokenverify.h"

int 
main(int argc, char **argv) 
{
    char idtoken[4000];
    char nonce[100];
    char *emsg;
    char *uid;
    int ret = -1;

    ret = setOAuthConfigC(getenv("OAUTH2_CLIENT_ID"), "https://accounts.google.com", &emsg);
    if (ret != 0) {
        printf("Failed to set auth config, emsg=%s\n", emsg);
        free(emsg);
        return ret;
    }

    printf("\nEnter ID token: ");
    scanf("%s", idtoken);

    printf("\nEnter nonce: ");
    scanf("%s", nonce);

    ret = verifyTokenC(idtoken, nonce, &uid, &emsg);
    if (ret != 0) {
        printf("Failed to set auth config, emsg=%s\n", emsg);
         free(emsg);
        return ret;
    }

    printf("userid verified=%s\n", uid);
    free(uid);

    return 0;
}