package main

import (
	"github.com/mitranim/gg"
)

var steamAppUrlBase = urlParse(`https://store.steampowered.com/app/`)
var steamAppImgUrlBase = urlParse(`https://cdn.cloudflare.steamstatic.com/steam/apps/`)

const steamAppImgUrlSuffix = `header.jpg`

type SteamApp struct {
	Appid uint32 `json:"appid"`
	Name  string `json:"name"`
}

func (self SteamApp) Game() (out Game) {
	id := gg.String(self.Appid)
	out.Name = self.Name
	out.Link = steamAppUrlBase.AddPath(id)
	out.Img = steamAppImgUrlBase.AddPath(id, steamAppImgUrlSuffix)
	return
}

func ReadSteamApps(path string) []SteamApp {
	return gg.JsonDecodeFileTo[[]SteamApp](path)
}

func ReadSteamGames(path string) []Game {
	return gg.Map(ReadSteamApps(path), SteamApp.Game)
}
