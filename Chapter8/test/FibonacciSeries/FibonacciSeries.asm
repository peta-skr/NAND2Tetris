@256
D=A
@SP
M=D
@300
D=A
@LCL
M=D
@400
D=A
@ARG
M=D
@3000
D=A
@THIS
M=D
@3010
D=A
@THAT
M=D
@ARG
D=M
@1
D=D+A
A=D
D=M
@SP
A=M
M=D
@ARG
D=M
@SP
M=M+1
@SP
M=M-1
@SP
A=M
D=M
M=0
@THAT
M=D
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@0
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
M=0
@R13
A=M
M=D
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@1
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
M=0
@R13
A=M
M=D
@ARG
D=M
@0
D=D+A
A=D
D=M
@SP
A=M
M=D
@ARG
D=M
@SP
M=M+1
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
@ARG
D=M
@0
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
M=0
@R13
A=M
M=D
@R13
M=0
(LOOP)
@ARG
D=M
@0
D=D+A
A=D
D=M
@SP
A=M
M=D
@ARG
D=M
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@COMPUTE_ELEMENT
D;JNE
@END
0;JMP
(COMPUTE_ELEMENT)
@THAT
D=M
@0
D=D+A
A=D
D=M
@SP
A=M
M=D
@THAT
D=M
@SP
M=M+1
@THAT
D=M
@1
D=D+A
A=D
D=M
@SP
A=M
M=D
@THAT
D=M
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
@THAT
D=M
@2
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
M=0
@R13
A=M
M=D
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
@SP
M=M-1
@SP
A=M
D=M
M=0
@THAT
M=D
@ARG
D=M
@0
D=D+A
A=D
D=M
@SP
A=M
M=D
@ARG
D=M
@SP
M=M+1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
@ARG
D=M
@0
D=D+A
@R13
M=D
@SP
M=M-1
@SP
A=M
D=M
M=0
@R13
A=M
M=D
@R13
M=0
@LOOP
0;JMP
(END)
