#!/usr/bin/env bash
./.github/scripts/000.create_identity.sh && \
./.github/scripts/100.fix_version.sh && \
./.github/scripts/200.run_tests.sh && \
./.github/scripts/300.create_zip.sh && \
./.github/scripts/400.create_tag.sh && \
./.github/scripts/500.save_changes.sh