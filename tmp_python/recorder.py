# 此原型工作基本正常 但仍需调整
# F: 添加机型适配,自动分辨率与横屏坐标等功能
import time
import subprocess
import sys

def get_xy():
    cmd = r'adb shell getevent'
    # adb触屏信息 adb shell getevent -p '0035 0036'
    xmin = 0
    xmax = 5759
    ymin = 0
    ymax = 12863
    # 屏幕分辨率相关
    resW = 1080
    resH = 2412
    # 是否横屏充电口向左1+10P
    # landscape = False
    landscape = True
    w = 0
    h = 0
    start = time.time()
    last = 0
    f = open(f'{sys.argv[1]}.dat', 'w')
    try:
        p1=subprocess.Popen(cmd,shell=True,stdout=subprocess.PIPE)
        for line in p1.stdout:
            line = line.decode(encoding="utf-8", errors="ignore")
            line = line.strip()

            # if ' 014a 00000001' in line:
            #     print('DOWN')

            if ' 0035 ' in line:
                e = line.split(" ")
                w = e[3]
                w = int(w, 16)
                w = (w - xmin) * resW/(xmax-xmin)
                
            if  ' 0036 ' in line:
                e = line.split(" ")
                h = e[3]
                h = int(h, 16)
                h = (h - ymin) * resH/(ymax-ymin)


            
            if ' 014a 00000000' in line:
                # print('UP')
                print('======================')
                f.write('======================\n')
            if landscape == False:
                elapsed = time.time() - start - last
                if elapsed <= 0.5:
                    continue
                print(int(elapsed))
                last = time.time() - start
                if int(w) == 0 | int(h) == 0:
                    continue
                print(f'X: {round(w, 1)} , Y: {round(h, 1)} ')
                f.write(f'@ {int(elapsed)}\n')
                f.write(f'# {round(w, 1)} {round(h, 1)}\n')
            if landscape == True:
                x = abs(h - resH)
                y = abs(w)
                if int(x) == 0 | int(y) == 0:
                    continue
                elapsed = time.time() - start - last
                if elapsed <= 0.5:
                    continue
                print(int(elapsed))
                last = time.time() - start
                print(f'X: {round(x, 1)} , Y: {round(y, 1)} ')
                f.write(f'@ {int(elapsed)}\n')
                f.write(f'# {round(x, 1)} {round(y, 1)}\n')
            
        p1.wait()
        
    except KeyboardInterrupt:
        f.close
        sys.exit()
    except Exception as e:
        f.close
        print(e)

get_xy()
# print(f'X: {size[0]} , Y: {size[1]} ') 