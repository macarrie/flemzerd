// TODO
// [] Pouvoir vérifier la validité de la conf (bonne clés pour les notifiers et les providers par exemple)
//      [] Bonne clés pour pushbullet, tvdb
//      [] Au moins un show de configuré
// [] Faire une boucle infinie pour regarder regulierement les nouveaux episodes et notifier si nouveaux resultats (poll_interval configurable)
//      [] Charger poll_interval dans la conf
//      [] Utiliser poll_interval comme intervalle de boucle
//      [] Dans la boucle, checker les derniers episodes. Si on a deja envoyé une notif pour cet episode, ne rien faire. Sinon, envoyer une notif comme quoi l'episode est dispo
// [] Voir comment marche jackett et si ya pas moyen de rechercher des torrents avec
// [] Si on peut les recuperer, faire une interface Downloaders et lancer un telechargement dessus (démon transmission pour commencer)
// [] Pouvoir configurer un dossier destination
// [] Faire des hook en fin de téléchargement (maj kodi, création d'un dossier dans un dossier destination prédéfini et copier le fichier dedans)
// [] Transformer tvdb en provider
// [] Faire une interface provider pour charger choisir un provider (et un seul ? sinon c'est chiant a gerer et ca sert a rien)
// [] Tests unitaires
// [] Doc ?

// TODO List de riche
// [] Retention fichier des derniers episodes notifiés. Si le demon tombe, pas de nouvelle notif pour rien
// [] Script systemd pour transformer l'exec en démon et le gerer avec systemctl
// [] Ajouter un autre notifier (mail par exemple)
// [] Interface web pour voir vite fait ce qu'il se passe: nouveaux episodes, shows surveillés et telechargements en cours

package main

import (
	"fmt"
	//"os"
	flag "github.com/ogier/pflag"
	"time"
	//"io/ioutil"
	"flemzerd/configuration"
	log "flemzerd/logging"
	"flemzerd/notifier"
	"flemzerd/pushbullet"
	"flemzerd/tvdb"
)

var config configuration.Configuration
var notificationsRetention []int

func initProviders(config configuration.Configuration) {
	log.Info("Initializing Providers")
}

func initDownloaders(config configuration.Configuration) {
	log.Info("Initializing Downloaders")
}

func initNotifiers(config configuration.Configuration) {
	log.Info("Initializing Notifiers")
	//fmt.Printf("%#v\n", config)
	for name, notifierObject := range config.Notifiers {
		//fmt.Printf("Name: %v, Notifier: %#v\n", name, notifierObject)
		switch name {
		case "pushbullet":
			pushbulletNotifier := pushbullet.New(map[string]string{"AccessToken": notifierObject["accesstoken"]})
			notifier.AddNotifier(pushbulletNotifier)

			log.WithFields(log.Fields{
				"notifier": pushbulletNotifier,
			}).Info("Notifier added to list of notifiers")
		}
	}
}

func NotifyRecentEpisode(show tvdb.Show, episode tvdb.Episode) {
	episodeLogString := fmt.Sprintf("S%03dE%03d: %v", episode.AiredSeason, episode.AiredEpisodeNumber, episode.EpisodeName)

	alreadyNotified := false
	var cleanedRetention []int

	for episodeId := range notificationsRetention {
		airDate, err := time.Parse("2006-01-02", episode.FirstAired)
		if err != nil {
			continue
		}

		if airDate.Before(time.Now()) && airDate.After(time.Now().AddDate(0, 0, -14)) {
			cleanedRetention = append(cleanedRetention, episode.Id)
		}

		if episode.Id == episodeId {
			alreadyNotified = true

			break
		}
	}
	notificationsRetention = cleanedRetention

	if alreadyNotified {
		log.WithFields(log.Fields{
			"show":    show.SeriesName,
			"episode": episodeLogString,
		}).Info("Notifications already sent for episode. Nothing to do")

		return
	}

	notificationTitle := fmt.Sprintf("%v: New episode aired (S%03dE%03d)", show.SeriesName, episode.AiredSeason, episode.AiredEpisodeNumber)
	notificationContent := fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", episode.FirstAired, show.SeriesName, episode.AiredSeason, episode.AiredEpisodeNumber, episode.EpisodeName)

	log.Debug(notificationTitle)
	log.Debug(notificationContent)
	log.Debug("Notification sent debug")
	//err := notifier.SendNotification(notificationTitle, notificationContent)
	//if err != nil {
	//log.Warning("Failed to send all notifications")
	//} else {
	notificationsRetention = append(notificationsRetention, episode.Id)
	//}

}

func main() {
	debugMode := flag.BoolP("debug", "d", false, "Start in debug mode")
	configFilePath := flag.StringP("config", "c", "", "Configuration file path to use")

	flag.Parse()

	if *debugMode {
		log.Setup(true)
	} else {
		log.Setup(false)
	}

	if *configFilePath != "" {
		// TODO
		log.Info("Loading provided configuration file")
		configuration.UseFile(*configFilePath)
	}

	config, err := configuration.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
	}

	initNotifiers(config)
	initProviders(config)
	initDownloaders(config)

	if !tvdb.Authenticate(config.Providers["tvdb"]["apikey"], config.Providers["tvdb"]["username"], config.Providers["tvdb"]["userkey"]) {
		log.Fatal("Unable to get tvdb API token")
	}

	// Load configuration objects
	var showObjects []tvdb.Show

	for _, show := range config.Shows {
		show, err := tvdb.FindShow(show)
		if err != nil {
			log.Warning(err)
		} else {
			showObjects = append(showObjects, show)
		}
		log.Debug(show)
	}
	log.Debug("Shows objects loaded from config: ", showObjects)

	log.Debug("Starting polling loop")
	loopTicker := time.NewTicker(15 * time.Second)
	//stopMainLoop := make(chan struct{})
	for {
		log.Debug("========== Polling loop start ==========")

		for _, show := range showObjects {
			recentEpisode, err := tvdb.FindRecentlyAiredEpisodeForShow(show)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Warning("No recent episode found for show")
				continue
			}
			// Send notification
			NotifyRecentEpisode(show, recentEpisode)
		}

		log.Debug("========== Polling loop end ==========\n")
		<-loopTicker.C
	}

	log.Debug("Polling loop terminated. Shutting down daemon")
}
