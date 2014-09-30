###Game Database testing

Size: 100 (1m)

###Individual results
Each block is stored in its own key
b:{xpos}:{ypos}:{zpos}

    [Sep 30 14:30:25] Info: Time taken setting 1000000 blocks: 26.873166717s, Blocks/s: 37211.84
    [Sep 30 14:30:25] Info: Starting get test
    [Sep 30 14:30:51] Info: Time taken Getting 1000000 blocks: 25.639977526s, Blocks/s: 39001.59

Memory: 85m

###Chunked results
Blocks are stored in chunks (redis list)
c:{xpos}:{ypos}

    [Sep 30 14:47:01] Info: Time taken Setting 1000000 blocks: 1.333832957s, Blocks/s: 749719.07
    [Sep 30 14:47:01] Info: Starting Get test
    [Sep 30 14:47:03] Info: Time taken Getting 1000000 blocks: 1.31965752s, Blocks/s: 757772.37

Memory 46m