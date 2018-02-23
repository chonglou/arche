cmake_minimum_required (VERSION 3.10)
project (arche)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if(CMAKE_BUILD_TYPE STREQUAL "Release")
    set(CMAKE_CXX_LINK_FLAGS "${CMAKE_CXX_LINK_FLAGS} -s")
endif()

AUX_SOURCE_DIRECTORY(src DIR_SRCS)

ADD_EXECUTABLE(arche ${DIR_SRCS})
TARGET_LINK_LIBRARIES(
  arche
  pqxx
  PocoCrypto
  PocoData
  PocoDataODBC
  PocoFoundation
  PocoJSON
  PocoNet
  PocoNetSSL
  PocoRedis
  PocoUtil
  PocoXML
  PocoZip
)
