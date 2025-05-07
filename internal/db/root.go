package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/redjax/go-mygithub/internal/domain/Github"
)

// Model for Github repository owner
func ConvertOwnerToModel(owner Github.RepositoryOwner) Github.RepositoryOwnerModel {
	return Github.RepositoryOwnerModel{
		ID:                owner.ID,
		Login:             owner.Login,
		NodeID:            owner.NodeID,
		AvatarURL:         owner.AvatarURL,
		GravatarID:        owner.GravatarID,
		URL:               owner.URL,
		HTMLURL:           owner.HTMLURL,
		FollowersURL:      owner.FollowersURL,
		FollowingURL:      owner.FollowingURL,
		GistsURL:          owner.GistsURL,
		StarredURL:        owner.StarredURL,
		SubscriptionsURL:  owner.SubscriptionsURL,
		OrganizationsURL:  owner.OrganizationsURL,
		ReposURL:          owner.ReposURL,
		EventsURL:         owner.EventsURL,
		ReceivedEventsURL: owner.ReceivedEventsURL,
		Type:              owner.Type,
		UserViewType:      owner.UserViewType,
		SiteAdmin:         owner.SiteAdmin,
	}
}

// Model for Github repository license
func ConvertLicenseToModel(license *Github.RepositoryLicense) *Github.RepositoryLicenseModel {
	if license == nil {
		return nil
	}
	return &Github.RepositoryLicenseModel{
		Key:    toNullString(license.Key),
		Name:   toNullString(license.Name),
		SPDXID: toNullString(license.SPDXID),
		URL:    toNullString(license.URL),
		NodeID: toNullString(license.NodeID),
	}
}

// Model for Github repository permissions
func ConvertPermissionsToModel(perm *Github.RepositoryPermissions) *Github.RepositoryPermissionsModel {
	if perm == nil {
		return nil
	}
	return &Github.RepositoryPermissionsModel{
		Admin:    toNullBool(perm.Admin),
		Maintain: toNullBool(perm.Maintain),
		Push:     toNullBool(perm.Push),
		Triage:   toNullBool(perm.Triage),
		Pull:     toNullBool(perm.Pull),
	}
}

// Model for Github repository
func ConvertRepositoryToModel(repo Github.Repository) Github.RepositoryModel {
	ownerModel := ConvertOwnerToModel(repo.Owner)
	licenseModel := ConvertLicenseToModel(repo.License)
	permModel := ConvertPermissionsToModel(repo.Permissions)

	homepage := toNullString(repo.Homepage)
	description := toNullString(repo.Description)
	language := toNullString(repo.Language)
	mirrorURL := toNullString(repo.MirrorURL)

	topicsJSON, _ := json.Marshal(repo.Topics) // marshal topics slice to JSON

	return Github.RepositoryModel{
		ID:                       repo.ID,
		NodeID:                   repo.NodeID,
		Name:                     repo.Name,
		FullName:                 repo.FullName,
		Private:                  repo.Private,
		OwnerID:                  ownerModel.ID,
		Owner:                    &ownerModel,
		HTMLURL:                  repo.HTMLURL,
		Description:              description,
		Fork:                     repo.Fork,
		URL:                      repo.URL,
		ForksURL:                 repo.ForksURL,
		KeysURL:                  repo.KeysURL,
		CollaboratorsURL:         repo.CollaboratorsURL,
		TeamsURL:                 repo.TeamsURL,
		HooksURL:                 repo.HooksURL,
		IssueEventsURL:           repo.IssueEventsURL,
		EventsURL:                repo.EventsURL,
		AssigneesURL:             repo.AssigneesURL,
		BranchesURL:              repo.BranchesURL,
		TagsURL:                  repo.TagsURL,
		BlobsURL:                 repo.BlobsURL,
		GitTagsURL:               repo.GitTagsURL,
		GitRefsURL:               repo.GitRefsURL,
		TreesURL:                 repo.TreesURL,
		StatusesURL:              repo.StatusesURL,
		LanguagesURL:             repo.LanguagesURL,
		StargazersURL:            repo.StargazersURL,
		ContributorsURL:          repo.ContributorsURL,
		SubscribersURL:           repo.SubscribersURL,
		SubscriptionURL:          repo.SubscriptionURL,
		CommitsURL:               repo.CommitsURL,
		GitCommitsURL:            repo.GitCommitsURL,
		CommentsURL:              repo.CommentsURL,
		IssueCommentURL:          repo.IssueCommentURL,
		ContentsURL:              repo.ContentsURL,
		CompareURL:               repo.CompareURL,
		MergesURL:                repo.MergesURL,
		ArchiveURL:               repo.ArchiveURL,
		DownloadsURL:             repo.DownloadsURL,
		IssuesURL:                repo.IssuesURL,
		PullsURL:                 repo.PullsURL,
		MilestonesURL:            repo.MilestonesURL,
		NotificationsURL:         repo.NotificationsURL,
		LabelsURL:                repo.LabelsURL,
		ReleasesURL:              repo.ReleasesURL,
		DeploymentsURL:           repo.DeploymentsURL,
		CreatedAt:                repo.CreatedAt,
		UpdatedAt:                repo.UpdatedAt,
		PushedAt:                 repo.PushedAt,
		GitURL:                   repo.GitURL,
		SshURL:                   repo.SshURL,
		CloneURL:                 repo.CloneURL,
		SvnURL:                   repo.SvnURL,
		Homepage:                 homepage,
		Size:                     repo.Size,
		StargazersCount:          repo.StargazersCount,
		WatchersCount:            repo.WatchersCount,
		Language:                 language,
		HasIssues:                repo.HasIssues,
		HasProjects:              repo.HasProjects,
		HasDownloads:             repo.HasDownloads,
		HasWiki:                  repo.HasWiki,
		HasPages:                 repo.HasPages,
		HasDiscussions:           repo.HasDiscussions,
		ForksCount:               repo.ForksCount,
		MirrorURL:                mirrorURL,
		Archived:                 repo.Archived,
		Disabled:                 repo.Disabled,
		OpenIssuesCount:          repo.OpenIssuesCount,
		LicenseID:                nil, // set after License saved
		License:                  licenseModel,
		AllowForking:             repo.AllowForking,
		IsTemplate:               repo.IsTemplate,
		WebCommitSignoffRequired: repo.WebCommitSignoffRequired,
		Topics:                   topicsJSON,
		Visibility:               repo.Visibility,
		Forks:                    repo.Forks,
		OpenIssues:               repo.OpenIssues,
		Watchers:                 repo.Watchers,
		DefaultBranch:            repo.DefaultBranch,
		PermissionsID:            nil, // set after Permissions saved
		Permissions:              permModel,
	}
}

// Helper to convert *string to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// Helper to convert *bool to sql.NullBool
func toNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

// Save Repository models to database
func SaveRepositories(db *gorm.DB, repos []Github.Repository) error {
	for i, repo := range repos {
		// Convert schemas to model
		model := ConvertRepositoryToModel(repo)

		// Save owner if present in model & not in database
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(model.Owner).Error; err != nil {
			return fmt.Errorf("repo %d (owner): %w", i+1, err)
		}
		model.OwnerID = model.Owner.ID

		// Save license if present in model & not in database
		if model.License != nil {
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(model.License).Error; err != nil {
				return fmt.Errorf("repo %d (license): %w", i+1, err)
			}
			model.LicenseID = &model.License.ID
		}

		// Save permissions if present in model & not in database
		if model.Permissions != nil {
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(model.Permissions).Error; err != nil {
				return fmt.Errorf("repo %d (permissions): %w", i+1, err)
			}
			model.PermissionsID = &model.Permissions.ID
		}

		// Save repo if it doesn't exist
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model).Error; err != nil {
			return fmt.Errorf("repo %d (main): %w", i+1, err)
		}

		// Log every 100 repos
		if (i+1)%100 == 0 || i == len(repos)-1 {
			fmt.Printf("  Saved %d/%d repositories to DB...\n", i+1, len(repos))
		}
	}

	return nil
}

// Initialize the database
func InitDB() (*gorm.DB, error) {
	// Create database connection
	db, err := gorm.Open(sqlite.Open("mygithub.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Do migrations
	err = db.AutoMigrate(
		&Github.RepositoryModel{},
		&Github.RepositoryOwnerModel{},
		&Github.RepositoryLicenseModel{},
		&Github.RepositoryPermissionsModel{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
