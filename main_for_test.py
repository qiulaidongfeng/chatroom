from channel.channel import *
import time


def main():
    start = time.time()
    c = Channel()
    print(time.time() - start)

    start = time.time()
    c.CreateRoom("test")
    print(time.time() - start)

    start = time.time()
    c.SendMessage("test", "k")
    print(time.time() - start)

    start = time.time()
    c.waitMessage()
    print(time.time() - start)

    start = time.time()
    msg = c.GetHistory("test")
    print(time.time() - start)

    start = time.time()
    c.ExitRoom("test")
    print(time.time() - start)

    assert msg != None
    assert msg[0] == "k"
    return


if __name__ == "__main__":
    main()
