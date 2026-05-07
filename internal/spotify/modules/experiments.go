package modules

var disableExperiments = []string{
	"enableInAppMessaging",
	"enableContentInformationMessage",
	"enablePickAndShuffle",
	"enableDesktopMusicLeavebehinds",
	"enableHptoLocationRefactor",
	"enableUserFraudSignals",
	"enableUserFraudVerificationRequest",
	"enableUserFraudVerification",
	"enableUserFraudCspViolation",
	"enableEsperantoMigration",
	"enableBillboardEsperantoMigration",
	"enableEsperantoMigrationLeaderboard",
	"enableSponsoredPlaylistEsperantoMigration",
	"enableNewAdsNpv",
	"enableNewAdsNpvVideoTakeover",
	"enableNewAdsNpvColorExtraction",
	"enableAudiobookAdExclusivity",
	"enableNewAdsNpvNewVideoTakeoverSlot",
	"enableFraudLoadSignals",
	"enableGabitoAdEvent",
	"enableYourListeningUpsell",
	"podcastads-ads_npb",
	"podcastaudioplus-episode_entity",
	"podcastaudioplus-show_page",
	"AutoSeekToAdPosition",
	"enablePodcastSponsoredContent",
	"enableHomeAds",
	"enableLearningHomeCard",
	"enablePipImpressionLogging",
	"allowSwitchingBetweenHomeAdsAndHpto",
	"enableLimitedAdsLabelsOnPlaylistCards",
	"enableLimitedAdsLabelsOnSearch",
	"enable_ad_feedback_milestone_1",
	"enableLeaderboardCrossOriginIframe",
	"enableLyricsUpsell",
	"enableArtistNPVImpressions",
	"enableSpotlightImpressionLogging",
	"enableEnhanceLikedSongs",
	"enableEnhancePlaylistProd",
	"enableSurveyAds",
	"enableHomeImpressions",
	"enableSearchImpressions",
	"enableNewAdsNpvCanvasAds",
	"enableCanvasAds",
	"enableConnectedStateObserver",
	"enableEsperantoAdStateReportManager",
	"enableLeavebehindsMockData",
	"enableEmbeddedNpvAds",
	"enableEnhancedAdsClientDeconfliction",
	"enableAdCountUi",
	"enableSaxLeaderboardAds",
	"enableEmbeddedAdVisibilityLogging",
	"enableHpto",
	"enableSponsoredPlaylistV2",
	"enableSponsoredPlaylistV2ScrollCard",
	"enableEmbeddedAdsCarousel",
	"enableEmbeddedAdsFetchingOverCanvas",
	"enableSponsoredPlaylistMockEndpoint",
	"embeddedAdImpressionDoesNotIgnoreVisilibility",
	"enable_ad_feedback_in_stream",
	"enableHarmonyMediaPlaybackReporting",
	"enableDsa",
	"enableDsaAds",
	"enableDSASetting",
	"enableRightSidebarMerchFallback",
	"bypassApplyUpdateCheck",
	"enableEFlag",
	"enableDynamicColors",
}

var enableExperiments = []string{
	"enableHomeViaGraphQLV2",
	"enableBrowseViaPathfinder",
	"enableIgnoreInRecommendations",
	"enableEqualizer",
	"enableCarouselsOnHome",
	"enableAttackOnTitanEasterEgg",
	"enableAlbumReleaseAnniversaries",
	"enableYLXSidebar",
	"enableRightSidebar",
	"enableRightSidebarLyrics",
	"enableRightSidebarExtractedColors",
	"enableSilenceTrimmer",
	"enableSmallPlaybackSpeedIncrements",
	"enableShowFollowsSetting",
	"enableRightSidebarCredits",
	"enableWhatsNewFeed",
	"enableRightSidebarArtistEnhanced",
	"enableNewEntityHeaders",
	"enableReadAlongTranscripts",
	"enableRightSidebarTransitionAnimations",
	"enableYLXEnhancements",
	"enableConcertsInterested",
	"enableConcertsForThisIsPlaylist",
	"enableAlbumCoverArtModal",
	"enableSmartShuffle",
	"enableConcertsTicketPrice",
	"enableDynamicNormalizer",
	"enableHeBringsNpb",
	"enableAlbumPrerelease",
	"enableNpvAboutPodcast",
	"enableQueueOnRightPanel",
	"enableRecentlyPlayedShortcut",
	"enableAlignedCuration",
	"enableRemoteDownloads",
}

func buildExperimentsJS() string {
	disableJS := "["
	for i, exp := range disableExperiments {
		if i > 0 {
			disableJS += ","
		}
		disableJS += `"` + exp + `"`
	}
	disableJS += "]"

	enableJS := "["
	for i, exp := range enableExperiments {
		if i > 0 {
			enableJS += ","
		}
		enableJS += `"` + exp + `"`
	}
	enableJS += "]"

	return `
(function() {
	if (window.__spotiliteExpInstalled) return;
	window.__spotiliteExpInstalled = true;

	var disable = ` + disableJS + `;
	var enable = ` + enableJS + `;

	function setExpValue(values, name, value) {
		if (!values) return;
		if (values.set && typeof values.set === 'function') {
			values.set(name, value);
		} else if (values[name]) {
			values[name].value = value;
		}
	}

	function applyExperiments(values) {
		if (!values) return;
		disable.forEach(function(name) { setExpValue(values, name, false); });
		enable.forEach(function(name) { setExpValue(values, name, true); });
	}

	function tryHookRemoteConfig() {
		try {
			var sp = window.__spotify || window.spotify;
			if (!sp) return false;
			var config = null;
			if (sp.getRemoteConfigResolver && typeof sp.getRemoteConfigResolver === 'function') {
				config = sp.getRemoteConfigResolver();
			} else if (sp.getRemoteConfiguration && typeof sp.getRemoteConfiguration === 'function') {
				config = sp.getRemoteConfiguration();
			}
			if (config) {
				var values = config.values || config.activeProperties;
				if (values) {
					applyExperiments(values);
					(window.Spotx || (window.Spotx = {})).RemoteExp = values;
					return true;
				}
			}
		} catch(e) {}
		return false;
	}

	if (!tryHookRemoteConfig()) {
		var observer = new MutationObserver(function() {
			if (tryHookRemoteConfig()) {
				observer.disconnect();
			}
		});
		observer.observe(document.documentElement, { childList: true, subtree: true });
		setTimeout(function() { observer.disconnect(); }, 30000);
	}
})();
`
}

type ExperimentModule struct {
	BaseModule
}

func NewExperimentModule(enabled bool) *ExperimentModule {
	return &ExperimentModule{
		BaseModule: BaseModule{name: "experiments", enabled: enabled},
	}
}

func (m *ExperimentModule) JS() string { return buildExperimentsJS() }
