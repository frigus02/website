---
title: Export a client certificate from Firefox on Android without root
summary: Exporting a client certificate, which is installed in Firefox on an Android phone can be difficult. Especially when you need the private key and your phone is not rooted. Here's a way.
datetime: 2015-11-22T15:04:00+01:00
tags:
    - Firefox
    - Android
    - Client Certificate
    - ADB
---

When I set up an SSL certificate for this website I created an account on startssl.com.
I finished the account creation on my smartphone with Firefox. During this process a
client certificate was added to my browser and startssl.com suggested to export it, so
I can use it on other browsers and computers as well.

Now here comes the big question: how do you export a client certificate from Firefox on
Android?

Turns out it's possible (even without rooting your phone). Here are the steps, I had to
take. Note that I'm doing this in a Windows 10 PC, but it should work the same way on any
other OS.

1.  Install the Android SDK and USB drivers for your smartphone. You need to connect to
    it using ADB.

2.  Activate developer options and USB debugging on your smartphone.

3.  Connect your smartphone to your computer using USB cable and make a backup of your
    Firefox app using ADB:

        adb backup org.mozilla.firefox

4.  Download the [Android Backup Extractor](https://sourceforge.net/projects/adbextractor/).
    It's a tool to convert ADB backup files to tar archives, which you can easily extract.
    Use it with your backup:

        java -jar .\android-backup-extractor-20151102-bin\abe.jar unpack .\backup.ab .\backup.tar

    Then extract the tar archive with any tool you want (e.g. 7-Zip).

5.  Now browse the backup and Locate your Firefox profile. It should be in a folder similar to this:

        apps\org.mozilla.firefox\f\mozilla\dx7e9l2j.default

    You should be able to find the files cert9.db and key4.db in it. These files contain the
    client certificate. In order to extract it you need to the pk12util from Mozilla's
    [NSS tools](https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS#Tools.2C_testing.2C_and_other_technical_details).

6.  Download and build the NSS tools. The required process is described on the Mozilla wiki in
    detail: [NSS sources building testing](https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS/NSS_Sources_Building_Testing).
    In short these are the steps, I followed:

    1.  Clone the repositories:

            mkdir nss-with-nspr
            cd nss-with-nspr
            hg clone https://hg.mozilla.org/projects/nspr
            hg clone https://hg.mozilla.org/projects/nss

    2.  Install MozillaBuildSetup-Latest.exe

    3.  Open start-shell-msvc2013.bat from mozilla build tools, navigate to the created nss-with-nspr\nss
        and folder and execute the command:

            make nss_build_all

    4.  Find pk12util.exe and certutil.exe in a folder similar to nss-with-nspr\dist\WIN954.0_DBG.OBJ\bin.

7.  This step is probably not necessary, when you correctly include the generated libraries in the path,
    but I was only able to execute pk12util, when my working directory was nss-with-nspr\dist\WIN954.0_DBG.OBJ\lib.

    Because of this I copied the files cert9.db and key4.db from my extracted Firefox profile to this
    lib folder.

8.  Then you can extract the certificate with the following commands:

        ..\bin\certutil.exe -K -d sql:.

    This will show all available certificates. Copy the certificate name from the last column and
    execute (enter any password you like, when asked):

        ..\bin\pk12util.exe -o cert.p12 -n "YOUR CERTIFICATE NAME" -d sql:.

9.  Congratulations. You should now have a file cert.p12 in your current working directory, which you
    can import in your Firefox on Desktop.
