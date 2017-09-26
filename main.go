package main

import (
    //"fmt"
    //"os"
    flag "github.com/ogier/pflag"
    //"io/ioutil"
    log "flemzerd/logging"
    "flemzerd/configuration"
    "flemzerd/tvdb"
    "flemzerd/notifier"
    "flemzerd/pushbullet"
)

var config configuration.Configuration

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

func main() {
    debugMode := flag.BoolP("debug", "d", false, "Start in debug mode")
    configFilePath := flag.StringP("config", "c", "", "Configuration file path to use")


    // TODO: Controler la validité de la conf

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
    //log.Info(config)

    initNotifiers(config)
    initProviders(config)
    initDownloaders(config)

    // TODO
    // [] Pouvoir vérifier la validité de la conf (bonne clés pour les notifiers et les providers par exemple)
    // [] Transformer tvdb en provider 
    // [] Faire une interface provider pour charger choisir un provider (et un seul ? sinon c'est chiant a gerer et ca sert a rien)
    // [] Faire une boucle infinie pour regarder regulierement les nouveaux episodes et notifier si nouveaux resultats (poll_interval configurable)
    // [] Envoyer une notif lorsqu'un nouvel episode est trouvé
    // [] Voir comment marche jackett et si ya pas moyen de rechercher des torrents avec
    // [] Si on peut les recuperer, faire une interface Downloaders et lancer un telechargement dessus (démon transmission pour commencer)
    // [] Pouvoir configurer un dossier destination
    // [] Faire des hook en fin de téléchargement (maj kodi, création d'un dossier dans un dossier destination prédéfini et copier le fichier dedans)
    // [] Tests unitaires
    // [] Doc ?

    // TODO List de riche
    // [] Script systemd pour transformer l'exec en démon et le gerer avec systemctl
    // [] Ajouter un autre notifier (mail par exemple)
    // [] Interface web pour voir vite fait ce qu'il se passe: nouveaux episodes, shows surveillés et telechargements en cours

    if !tvdb.Authenticate(config.Providers["tvdb"]["apikey"], config.Providers["tvdb"]["username"], config.Providers["tvdb"]["userkey"]) {
        log.Fatal("Unable to get tvdb API token")
    }

    for _, show := range config.Shows {
        nextEpisode, err := tvdb.FindNextEpisodeForShow(show)
        if err != nil {
            log.WithFields(log.Fields{
                "error": err,
            }).Warning("No next episode found for show")
        }
        log.Info(nextEpisode)
    }
}
