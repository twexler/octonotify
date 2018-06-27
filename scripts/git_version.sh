#!/usr/bin/env bash
[[ "${TRAVIS}" == "true" ]] && {
    [[ "${TRAVIS_TAG}" != "" ]] && { echo ${TRAVIS_TAG}; exit; }
    echo "${TRAVIS_COMMIT}" | cut -c -7
} || {
    echo $(git describe)
}