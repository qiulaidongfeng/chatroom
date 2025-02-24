from channel.channel import *


def test_Channel():
    c = Channel()
    c.CreateRoom("test")
    c.SendMessage("test", "k")
    c.waitMessage()
    msg = c.EnterRoom("test")
    c.ExitRoom("test")
    assert msg != None
    assert msg[0] == "k"
    return
