import hashlib
import time
import signal


def handler():
    raise Exception("Timeout")


msg = "nicolas458"
i = 0
msg + str(i)


debut = time.time()


counter = 0


msg = msg + str(i)
for j in range(3):
    while hashlib.sha256(msg.encode()).hexdigest()[:2] != "00":
        i += 1
        msg = msg + str(i)
    msg = msg + str(i)


fin = time.time()

print(hashlib.sha256(msg.encode()).hexdigest())
print(fin - debut)
print("ok c'est bon")
