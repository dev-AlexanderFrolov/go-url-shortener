# настройка связи гита с новым устройством

git config --global user.name "Твой Ник"
git config --global user.email "твой.email@example.com"

ssh-keygen -t ed25519 -C "твой.email@example.com"

cat ~/.ssh/id_ed25519.pub

git clone git@github.com:ТВОЙ_НИК/go-url-shortener.git