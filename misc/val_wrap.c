#include <stdio.h>
#include <string.h>
#include <libgen.h>

int 
main(int argc, char *argv[])
{
	char buf[1000];
	int i;
	char *short_cmd;

	short_cmd = basename(strdup(argv[0]));
	sprintf(buf, "valgrind --log-file=/tmp/%s_%d.val --suppressions=./valgrind.supp --leak-check=full %s.orig", short_cmd, getpid(), argv[0]);
	i = 1;
	while (i < argc) {
		strcat(buf, " ");
		strcat(buf, argv[i]);	
		i++;
	}

	system(buf);
}
