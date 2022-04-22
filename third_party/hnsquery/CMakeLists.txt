cmake_minimum_required(VERSION 3.20)
project(test C)

set(CMAKE_C_STANDARD 99)

if (MSVC)
    add_compile_options(/W4 /WX)
else()
    add_compile_options(-Wall -Wextra -pedantic -Werror)
endif()

include_directories(build/include)
include_directories(build/include/hsk)
include_directories(build/include/hsk/chacha20)
include_directories(build/include/hsk/poly1305)
include_directories(build/include/hsk/secp256k1)

add_executable(test test.c hns.c hns.h store.h store.c bio.h)
target_link_libraries(test ${CMAKE_SOURCE_DIR}/build/lib/libhsk.a ${CMAKE_SOURCE_DIR}/build/lib/libuv.a -ldl -lpthread)

