# you can use this script to generate "bundled" Go + bee binaries.
# aka. the interpreter with a fixed script, being your bee script.
# made in 5 minutes: report any bugs!!!

# import required libraries
import shutil
import os

script = input("script: ")
with open(script, "r") as f:
    script = f.read().replace("\n", "")
    f.close()

os.mkdir("./tmp")
for o in os.listdir():
    if o.endswith(".go") or o.endswith(".mod"):
        shutil.copyfile(o, f"tmp/{o}")

with open("./tmp/bee.go", "r") as f:
    data = f.read().split("\n")
    final = []
    for i in range(0, len(data)):
        if i > 37 and i < 43: # other var(...)
            continue
        elif i > 82 and i < 96: # check for files, etc.
            continue
        elif i > 102 and i < 108: # generateGlobal end
            continue
        elif i == 44: # program string
            final.append(f"    program = `{script}`")
        elif i > 109: # init(...)
            continue
        else:
            final.append(data[i])
    data = "\n".join(final).replace("\"flag\"", "")
    f.close()
with open("./tmp/bee.go", "w") as f:
    f.write(data)
    f.close()

# build
os.chdir("./tmp")
win = input("Compiling to Win? (y/n): ").lower().startswith("y")
os.system(f"{'env GOOS=windows GOARCH=amd64 ' if win else ''}go build -o ../bundled{'.exe' if os.system('cls') == 0 or win else ''} .")
os.chdir("..")

# cleanup
for o in os.listdir("./tmp"):
    os.remove(f"./tmp/{o}")
    pass
os.rmdir("./tmp")
