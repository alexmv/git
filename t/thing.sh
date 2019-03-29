#!/bin/sh

test_description='git status with file system watcher'

. ./test-lib.sh

clean_repo () {
	git reset --hard HEAD &&
	git clean -fd
}

dirty_repo () {
	: >untracked &&
	: >dir1/untracked &&
	: >dir2/untracked &&
	echo 1 >modified &&
	echo 2 >dir1/modified &&
	echo 3 >dir2/modified &&
	echo 4 >new &&
	echo 5 >dir1/new &&
	echo 6 >dir2/new
}

test_expect_success 'setup' '
	mkdir -p .git/hooks &&
	: >tracked &&
	: >modified &&
	mkdir dir1 &&
	: >dir1/tracked &&
	: >dir1/modified &&
	mkdir dir2 &&
	: >dir2/tracked &&
	: >dir2/modified &&
	git -c core.fsmonitor= add . &&
	git -c core.fsmonitor= commit -m initial &&
	git config core.fsmonitor .git/hooks/fsmonitor-test &&
	cat >.gitignore <<-\EOF
	.gitignore
	expect*
	actual*
	marker*
	EOF
'

# Ensure commands that call refresh_index() to move the index back in time
# properly invalidate the fsmonitor cache
test_expect_success 'refresh_index() invalidates fsmonitor cache' '
	write_script .git/hooks/fsmonitor-test<<-\EOF &&
	EOF
	clean_repo &&
	dirty_repo &&
	git add . &&
	git commit -m "to reset" &&
	git reset HEAD~1 &&
	git status >actual &&
	git -c core.fsmonitor= status >expect &&
	test_i18ncmp expect actual
'


test_done
