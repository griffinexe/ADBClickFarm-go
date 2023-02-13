# 此原型工作正常
# F: 添加行数计数
import os
import sys
import time

def tap(x, y):
    print(f'[ADB]: tap screen X:{x} Y:{y}')
    os.system(f'adb shell input tap {x} {y}')

def playScript(scName):
    f = open(f'{scName}')
    lines = f.readlines()
    delay = 0
    

    for line in lines:
        try:
            line = line.replace("\n","")
            print(f'[FILE]: {line}')
            if '=' in line:
                print(f'[EXEC]: end touch')
            if '@' in line:
                line = line.split(' ')
                # time.sleep(int(line[1]))
                delay = int(line[1])
            if '#' in line:
                line = line.split(' ')
                if line[1] == '0' or line[2] == '0':
                    print('[EXEC]: skipping invalid zero value')
                    continue
                tap(line[1], line[2])
                time.sleep(delay)
        except KeyboardInterrupt:
            sys.exit()


def play():
    execList = ['AK_farmChip_class1_universal.dat', 'AK_farmChip_class1_universal.dat', 'AK_farmChip_class1_universal.dat']
    for sc in execList:
        print(f'[EXEC]: playing script {sc}')
        playScript(sc)


play()