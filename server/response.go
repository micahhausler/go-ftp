package server

const (
	DataCnxAlreadyOpenStartXfr = "125 Data connection already open, starting transfer\r\n"
	TypeSetOk                  = "200 Type set ok\r\n"
	PortOk                     = "200 PORT ok\r\n"
	FeatResponse               = "211-Features:\r\n  FEAT\r\n  MDTM\r\n  PASV\r\n  SIZE\r\n  TYPE A;I\r\n211 End\r\n"
	SysType                    = "215 UNIX Type: L8\r\n"
	GoodbyeMsg                 = "221 Goodbye!"
	TxfrCompleteOk             = "226 Data transfer complete\r\n"
	CmdOk                      = "200 Command ok\r\n"
	EnteringPasvMode           = "227 Entering Passive Mode (%s)\r\n"
	PwdResponse                = "257 \"/\"\r\n"
	FtpServerReady             = "220 FTP Server Ready\r\n"
	UsrLoggedInProceed         = "230 User Logged In Proceed\r\n"
	UsrNameOkNeedPass          = "331 Username OK Need Pass\r\n"
	SyntaxErr                  = "500 Syntax Error\r\n"
	CmdNotImplmntd             = "502 Command not implemented\r\n"
	NotLoggedIn                = "530 Not Logged In\r\n"
	AuthFailure                = "530 Auth Failure\r\n"
	AuthFailureTryAgain        = "530 Please login with USER and PASS."
	AnonUserDenied             = "550 Anon User Denied\r\n"
)

/*
const (
	ServiceReadyInNMinutes  = 120
	DataCnxAlreadyOpenStartXfr  = 125
	FileStatusOkOpenDataCnx  = 150
	CmdOk  = 200.1
	TypeSetOk  = 200.2
	EnteringPortMode  = 200.3
	CmdNotImplmntdSuperfluous  = 202
	SysStatusOrHelpReply  = 211.1
	FeatOk  = 211.2
	DirStatus  = 212
	FileStatus  = 213
	HelpMsg  = 214
	NameSysType  = 215
	SvcReadyForNewUser  = 220.1
	WelcomeMsg  = 220.2
	SvcClosingCtrlCnx  = 221.1
	GoodbyeMsg  = 221.2
	DataCnxOpenNoXfrInProgress  = 225
	ClosingDataCnx  = 226.1
	TxfrCompleteOk  = 226.2
	EnteringPasvMode  = 227
	EnteringEpsvMode  = 229
	UsrLoggedInProceed  = 230.1
	GuestLoggedInProceed  = 230.2
	ReqFileActnCompletedOk  = 250
	PwdReply  = 257.1
	MkdReply  = 257.2
	UsrNameOkNeedPass  = 331.1
	GuestNameOkNeedEmail  = 331.2
	NeedAcctForLogin  = 332
	ReqFileActnPendingFurtherInfo  = 350
	SvcNotAvailClosingCtrlCnx  = 421.1
	TooManyConnections  = 421.2
	CantOpenDataCnx  = 425
	CnxClosedTxfrAborted  = 426
	ReqActnAbrtdFileUnavail  = 450
	ReqActnAbrtdLocalErr  = 451
	ReqActnAbrtdInsuffStorage  = 452
	SyntaxErr  = 500
	SyntaxErrInArgs  = 501
	CmdNotImplmntd  = 502.1
	OptsNotImplemented  = 502.2
	BadCmdSeq  = 503
	CmdNotImplmntdForParam  = 504
	NotLoggedIn  = 530.1
	AuthFailure  = 530.2
	NeedAcctForStor  = 532
	FileNotFound  = 550.1
	PermissionDenied  = 550.2
	AnonUserDenied  = 550.3
	IsNotADir  = 550.4
	ReqActnNotTaken  = 550.5
	FileExists  = 550.6
	IsADir  = 550.7
	PageTypeUnk  = 551
	ExceededStorageAlloc  = 552
	FilenameNotAllowed  = 553
)
*/
