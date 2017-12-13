#include "cache.h"

int cmd_main(int ac, const char **av)
{
	struct index_state *istate = &the_index;
	int i;

	setenv("GIT_FSMONITOR_TEST", "keep", 1);
	setup_git_directory();
	if (read_index(istate) < 0)
		die("unable to read index file");
	if (!istate->fsmonitor_last_update) {
		printf("no fsmonitor\n");
		return 0;
	}
	printf("fsmonitor last update %"PRIuMAX"\n", (uintmax_t)istate->fsmonitor_last_update);

	for (i = 0; i < istate->cache_nr; i++)
		printf((istate->cache[i]->ce_flags & CE_FSMONITOR_VALID) ? "+" : "-");

	return 0;
}
