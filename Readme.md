#arm build

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \  
 fyne package --os darwin \
 --icon assets/clip.png \
 --name "超吉貼" \
 --id com.yourdomain.superclip \
 --release

#amd build
arch -x86_64 /bin/zsh -l
go env GOARCH
uname -m

GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
fyne package --os darwin \
 --icon assets/clip.png \
 --name "超吉貼-intel" \
 --id com.yourdomain.superclip \
 --release

xattr -d com.apple.quarantine ~/Downloads/超吉貼.app

rm ~/Library/Application\ Support/ClipMini/auth.json
