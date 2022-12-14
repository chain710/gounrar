#ifndef _UNRAR_DLL_
#define _UNRAR_DLL_

#include "raros.hpp"
#include "wchar.h"
#include <inttypes.h>
#ifdef _WIN_ALL
#include <windows.h>
#endif
// for cgo to work: https://github.com/golang/go/wiki/cgo#struct-alignment-issues
#pragma pack(push, 8)

#define ERAR_SUCCESS             0
#define ERAR_END_ARCHIVE        10
#define ERAR_NO_MEMORY          11
#define ERAR_BAD_DATA           12
#define ERAR_BAD_ARCHIVE        13
#define ERAR_UNKNOWN_FORMAT     14
#define ERAR_EOPEN              15
#define ERAR_ECREATE            16
#define ERAR_ECLOSE             17
#define ERAR_EREAD              18
#define ERAR_EWRITE             19
#define ERAR_SMALL_BUF          20
#define ERAR_UNKNOWN            21
#define ERAR_MISSING_PASSWORD   22
#define ERAR_EREFERENCE         23
#define ERAR_BAD_PASSWORD       24
#define ERAR_ESEEK              25

#define RAR_OM_LIST              0
#define RAR_OM_EXTRACT           1
#define RAR_OM_LIST_INCSPLIT     2

#define RAR_SKIP              0
#define RAR_TEST              1
#define RAR_EXTRACT           2

#define RAR_VOL_ASK           0
#define RAR_VOL_NOTIFY        1

#define RAR_DLL_VERSION       8

#define RAR_HASH_NONE         0
#define RAR_HASH_CRC32        1
#define RAR_HASH_BLAKE2       2

#if defined(_UNIX) || defined(__APPLE__)
#define CALLBACK
#define PASCAL
#define LONG long
#define HANDLE void *
#define LPARAM long
#define UINT unsigned int
#endif

#define RHDF_SPLITBEFORE 0x01
#define RHDF_SPLITAFTER  0x02
#define RHDF_ENCRYPTED   0x04
#define RHDF_SOLID       0x10
#define RHDF_DIRECTORY   0x20


typedef struct
{
  char         ArcName[260];
  char         FileName[260];
  unsigned int Flags;
  unsigned int PackSize;
  unsigned int UnpSize;
  unsigned int HostOS;
  unsigned int FileCRC;
  unsigned int FileTime;
  unsigned int UnpVer;
  unsigned int Method;
  unsigned int FileAttr;
  char         *CmtBuf;
  unsigned int CmtBufSize;
  unsigned int CmtSize;
  unsigned int CmtState;
} RARHeaderData;


typedef struct
{
  char         ArcName[1024];
  wchar_t      ArcNameW[1024];
  char         FileName[1024];
  wchar_t      FileNameW[1024];
  unsigned int Flags;
  unsigned int PackSize;
  unsigned int PackSizeHigh;
  unsigned int UnpSize;
  unsigned int UnpSizeHigh;
  unsigned int HostOS;
  unsigned int FileCRC;
  unsigned int FileTime;
  unsigned int UnpVer;
  unsigned int Method;
  unsigned int FileAttr;
  char         *CmtBuf;
  unsigned int CmtBufSize;
  unsigned int CmtSize;
  unsigned int CmtState;
  unsigned int DictSize;
  unsigned int HashType;
  char         Hash[32];
  unsigned int RedirType;
  wchar_t      *RedirName;
  unsigned int RedirNameSize;
  unsigned int DirTarget;
  unsigned int MtimeLow;
  unsigned int MtimeHigh;
  unsigned int CtimeLow;
  unsigned int CtimeHigh;
  unsigned int AtimeLow;
  unsigned int AtimeHigh;
  unsigned int Type; // header type
  uint64_t MtimeUnix;
  uint64_t CtimeUnix;
  uint64_t AtimeUnix;
  int64_t BlockPos;
} RARHeaderDataEx;


typedef struct
{
  char         *ArcName;
  unsigned int OpenMode;
  unsigned int OpenResult;
  char         *CmtBuf;
  unsigned int CmtBufSize;
  unsigned int CmtSize;
  unsigned int CmtState;
} RAROpenArchiveData;

typedef int (CALLBACK *UNRARCALLBACK)(UINT msg,uintptr_t UserData,uintptr_t P1,uintptr_t P2);

#define ROADF_VOLUME       0x0001
#define ROADF_COMMENT      0x0002
#define ROADF_LOCK         0x0004
#define ROADF_SOLID        0x0008
#define ROADF_NEWNUMBERING 0x0010
#define ROADF_SIGNED       0x0020
#define ROADF_RECOVERY     0x0040
#define ROADF_ENCHEADERS   0x0080
#define ROADF_FIRSTVOLUME  0x0100

#define ROADOF_KEEPBROKEN  0x0001

typedef struct
{
  char         *ArcName;
  wchar_t      *ArcNameW;
  unsigned int  OpenMode;
  unsigned int  OpenResult;
  char         *CmtBuf;
  unsigned int  CmtBufSize;
  unsigned int  CmtSize;
  unsigned int  CmtState;
  unsigned int  Flags;
  UNRARCALLBACK Callback;
  uintptr_t        UserData;
  unsigned int  OpFlags;
  wchar_t      *CmtBufW;
  unsigned int  Reserved[25];
} RAROpenArchiveDataEx;

enum UNRARCALLBACK_MESSAGES {
  UCM_CHANGEVOLUME,UCM_PROCESSDATA,UCM_NEEDPASSWORD,UCM_CHANGEVOLUMEW,
  UCM_NEEDPASSWORDW
};

typedef int (PASCAL *CHANGEVOLPROC)(char *ArcName,int Mode);
typedef int (PASCAL *PROCESSDATAPROC)(unsigned char *Addr,int Size);

#ifdef __cplusplus
extern "C" {
#endif

void* PASCAL RAROpenArchive(RAROpenArchiveData *ArchiveData);
void* PASCAL RAROpenArchiveEx(RAROpenArchiveDataEx *ArchiveData);
int    PASCAL RARCloseArchive(void* hArcData);
int    PASCAL RARReadHeader(void* hArcData, RARHeaderData *HeaderData);
int    PASCAL RARReadHeaderEx(void* hArcData, RARHeaderDataEx *HeaderData);
int    PASCAL RARProcessFile(void* hArcData,int Operation,char *DestPath,char *DestName);
int    PASCAL RARProcessFileW(void* hArcData,int Operation,wchar_t *DestPath,wchar_t *DestName);
void   PASCAL RARSetCallback(void* hArcData,UNRARCALLBACK Callback,uintptr_t UserData);
void   PASCAL RARSetChangeVolProc(void* hArcData,CHANGEVOLPROC ChangeVolProc);
void   PASCAL RARSetProcessDataProc(void* hArcData,PROCESSDATAPROC ProcessDataProc);
void   PASCAL RARSetPassword(void* hArcData,char *Password);
int    PASCAL RARGetDllVersion();
int    PASCAL RARSeek(void* hArcData, int64_t blockPos);
//int    PASCAL RARPExtractFile(HANDLE hArcData,int Operation,char *DestPath,char *DestName);
#ifdef __cplusplus
}
#endif

#pragma pack(pop)

#endif
