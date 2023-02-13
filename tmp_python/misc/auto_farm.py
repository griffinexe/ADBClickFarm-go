import os
import time
import datetime

# os.system('adb')
# os.system('adb shell input tap 968.2 713.2')

# time.sleep(5)
# os.system('adb shell input tap 925.2 42.2')

def tap(x, y):
    print(f'tap screen X:{x} Y:{y}')
    os.system(f'adb shell input tap {x} {y}')

def Dtap(x, y, sec):
    print(f'tap screen X:{x} Y:{y} and delay {sec}S')
    os.system(f'adb shell input tap {x} {y}')
    time.sleep(sec)

def swipe(Ax, Ay, Bx, By, msec):
    print(f'swipe screen from X1:{Ax} Y1:{Ay} to X2:{Bx} Y2:{By} using {msec} ms')
    os.system(f'adb shell input swipe {Ax} {Ay} {Bx} {By} {msec}')

# def dragDrop(Ax, Ay, Bx, By, msec):
#     print(f'swipe screen from X1:{Ax} Y1:{Ay} to X2:{Bx} Y2:{By} using {msec} ms')
#     os.system(f'adb shell input swipe {Ax} {Ay} {Bx} {By} {msec}')

# ark VM
# Dtap(968.2, 713.2, 5)
# Dtap(925.2, 42.2, 3)
# Dtap(894.2, 722.2, 2)
# Dtap(894.2, 722.2, 8)
# Dtap(894.2, 722.2, 1)
# Dtap(40.2, 36.2, 2)
# Dtap(662.2, 723.2, 1)
# Dtap(735.2, 479.2, 8)
# Dtap(928.2, 595.2, 1)
# Dtap(928.2, 595.2, 1)

# AL device
# Get weekday
dayWeek = datetime.datetime.today().weekday()
# 复位
def reset():
    Dtap(3110.5, 72.0, 1)
    Dtap(2875, 728, 1)
# 返回
def goBack():
    Dtap(128, 138, 1)
# LOGIN
Dtap(2955, 1097, 2)
Dtap(2875, 728, 1)
# 演习近平
Dtap(2450, 1256, 1)
# 第一个对手*10
for x in range(10):
    Dtap(748, 568, 1)
    Dtap(1600, 1144, 1)
    Dtap(2961, 1296, 80)
    Dtap(2961, 1296, 1)
    Dtap(2961, 1296, 1)
    Dtap(2504, 1326, 1)
reset()
# 挑战
Dtap(1786, 1245, 1)
Dtap(1657, 771, 1)
Dtap(1261, 417, 1)
#TODO
Dtap(1888, 501, 1)
# Dtap(1888, 501, 1)
Dtap(2961, 1296, 1)
goBack()
if dayWeek in (2, 5, 6):
    Dtap(1090, 901, 1)
    Dtap(1261, 417, 1)
    Dtap(1888, 501, 1)
    Dtap(2961, 1296, 1)
    goBack()
if dayWeek in (1, 4, 6):
    Dtap(620, 932, 1)
    Dtap(1261, 417, 1)
    Dtap(1888, 501, 1)
    Dtap(2961, 1296, 1)
    goBack()
if dayWeek (0, 3, 6):
    Dtap(294, 944, 1)
    Dtap(1261, 417, 1)
    Dtap(1888, 501, 1)
    Dtap(2961, 1296, 1)
    goBack()
reset()
Dtap(973, 772, 1)
Dtap(680, 576, 1)
Dtap(2250, 1030, 60)
Dtap(2531, 1232)
Do...

# 侧菜单
Dtap(40, 350, 1)
# 收菜
Dtap(933, 175, 1)
Dtap(572, 161, 1)
Dtap(240, 191, 1)
# 委托收
Dtap(887, 538, 2)
Dtap(2961, 1296, 1)
Dtap(2961, 1296, 2)
Dtap(2961, 1296, 1)
Dtap(2961, 1296, 2)
Dtap(2961, 1296, 1)
Dtap(2961, 1296, 2)
Dtap(2961, 1296, 1)
Dtap(2961, 1296, 2)
# 委托发
Dtap(887, 538, 2)
