#include <stdio.h>
#include <stdlib.h>


int main(int argc, char *argv[])
{
    if (argc > 1) {
	for (int i=0; i < argc; i++) {       // This line is indented with tab
            printf("%d: %s\n", i, argv[i]);
        }
    }
    return 0;
}

