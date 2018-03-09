import argparse
import sys
import urllib.request
import os.path


def download(root, files):
    if not os.path.isdir(root):
        os.makedirs(root)
    for u, f in files.items():
        fn = os.path.join(root, f)
        if os.path.isfile(fn):
            print("File {} already exists".format(fn))
            continue
        print("Get {}".format(u))
        urllib.request.urlretrieve(u, filename=fn)


def run_download(name):
    if name == 'yinshun':
        print('Download 印顺全集')
        root = os.path.join('tmp', 'downloads', 'yinshun')
        files = {}
        for id in (list(range(1, 14 + 1)) + list(range(42, 44 + 1))):
            files["http://www.yinshun.org.tw/epub's%%20web/epub/y%02d.epub" %
                  id] = "y%02d.epub" % id
        download(root, files)
    elif name == 'dzdljh':
        print('Download 大智度论精华')
        root = os.path.join('tmp', 'downloads', 'T0332')
        files = {}
        for id in range(1, 2 + 1):
            files["http://ftp4.budaedu.org/ghosa4/C006/T0332/ref/T0332_%03d.pdf" %
                  id] = "T0332_%03d.pdf" % id
        for id in range(1, 28 + 1):
            files["http://ftp4.budaedu.org/ghosa4/C006/T0332/video-low/332%03dZ.mp4" %
                  id] = "332%03dZ.mp4" % id
            files["http://ftp4.budaedu.org/ghosa4/C006/T0332/audio-low/332%03dP.mp3" %
                  id] = "332%03dP.mp3" % id
        download(root, files)
    else:
        print('Unknown {}'.format(name))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='辅助工具。')
    parser.add_argument('-d', '--download',
                        choices=['yinshun', 'cbeta', 'dict', 'dzdljh'],
                        help='Download resources')
    args = parser.parse_args()
    if args.download:
        run_download(args.download)
    else:
        parser.print_help(sys.stderr)
