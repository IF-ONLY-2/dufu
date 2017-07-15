package main

const (
    PCI_COMMAND              uint16 = 0x04  /* 16 bits */
    PCI_COMMAND_IO                  = 0x1   /* Enable response in I/O space */
    PCI_COMMAND_MEMORY              = 0x2   /* Enable response in Memory space */
    PCI_COMMAND_MASTER              = 0x4   /* Enable bus mastering */
    PCI_COMMAND_SPECIAL             = 0x8   /* Enable response to special cycles */
    PCI_COMMAND_INVALIDATE          = 0x10  /* Use memory write and invalidate */
    PCI_COMMAND_VGA_PALETTE         = 0x20  /* Enable palette snooping */
    PCI_COMMAND_PARITY              = 0x40  /* Enable parity checking */
    PCI_COMMAND_WAIT                = 0x80  /* Enable address/data stepping */
    PCI_COMMAND_SERR                = 0x100 /* Enable SERR */
    PCI_COMMAND_FAST_BACK           = 0x200 /* Enable back-to-back writes */
    PCI_COMMAND_INTX_DISABLE uint16 = 0x400 /* INTx Emulation Disable */
)

// I/O+ Mem+ BusMaster+ SpecCycle- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisINTx+

func commandString(i uint16) string {
    var s string
    if i&PCI_COMMAND_IO != 0 {
        s += "I/O+ "
    } else {
        s += "I/O- "
    }
    if i&PCI_COMMAND_MEMORY != 0 {
        s += "Mem+ "
    } else {
        s += "Mem- "
    }
    if i&PCI_COMMAND_MASTER != 0 {
        s += "BusMaster+ "
    } else {
        s += "BusMaster- "
    }
    if i&PCI_COMMAND_SPECIAL != 0 {
        s += "SpecCycle+ "
    } else {
        s += "SpecCycle- "
    }
    if i&PCI_COMMAND_INVALIDATE != 0 {
        s += "MemWINV+ "
    } else {
        s += "MemWINV- "
    }

    if i&PCI_COMMAND_VGA_PALETTE != 0 {
        s += "VGASnoop+ "
    } else {
        s += "VGASnoop- "
    }
    if i&PCI_COMMAND_PARITY != 0 {
        s += "ParErr+ "
    } else {
        s += "ParErr- "
    }
    if i&PCI_COMMAND_WAIT != 0 {
        s += "Stepping+ "
    } else {
        s += "Stepping- "
    }
    if i&PCI_COMMAND_SERR != 0 {
        s += "SERR+ "
    } else {
        s += "SERR- "
    }
    if i&PCI_COMMAND_FAST_BACK != 0 {
        s += "FastB2B+ "
    } else {
        s += "FastB2B- "
    }
    if i&PCI_COMMAND_INTX_DISABLE != 0 {
        s += "DisINTx+ "
    } else {
        s += "DisINTx- "
    }
    s = s[:len(s)-1]
    return s
}
