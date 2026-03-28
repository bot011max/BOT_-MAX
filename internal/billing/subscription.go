package billing

type SubscriptionTier string

const (
    TierFree    SubscriptionTier = "free"
    TierPremium SubscriptionTier = "premium"
    TierFamily  SubscriptionTier = "family"
)

type Subscription struct {
    Tier      SubscriptionTier `json:"tier"`
    Price     float64          `json:"price"`
    Features  []string         `json:"features"`
    IsActive  bool             `json:"is_active"`
    ExpiresAt *time.Time       `json:"expires_at"`
}

var TierConfigs = map[SubscriptionTier]Subscription{
    TierFree: {
        Tier:  TierFree,
        Price: 0,
        Features: []string{
            "До 5 лекарств",
            "Базовые напоминания",
            "Email поддержка",
            "7 дней истории",
        },
    },
    TierPremium: {
        Tier:  TierPremium,
        Price: 299,
        Features: []string{
            "Неограниченное количество лекарств",
            "AI-анализ симптомов",
            "Приоритетная поддержка 24/7",
            "Экспорт отчетов в PDF",
            "Push-уведомления",
            "Безлимитная история",
        },
    },
    TierFamily: {
        Tier:  TierFamily,
        Price: 499,
        Features: []string{
            "Все функции Premium",
            "До 5 членов семьи",
            "Общий дашборд",
            "Приоритетная поддержка",
        },
    },
}

type SubscriptionManager struct {
    userSubscriptions map[string]*Subscription
}

func NewSubscriptionManager() *SubscriptionManager {
    return &SubscriptionManager{
        userSubscriptions: make(map[string]*Subscription),
    }
}

func (s *SubscriptionManager) GetUserTier(userID string) SubscriptionTier {
    if sub, exists := s.userSubscriptions[userID]; exists && sub.IsActive {
        return sub.Tier
    }
    return TierFree
}

func (s *SubscriptionManager) UpgradeToPremium(userID string) error {
    s.userSubscriptions[userID] = &Subscription{
        Tier:     TierPremium,
        Price:    299,
        Features: TierConfigs[TierPremium].Features,
        IsActive: true,
    }
    return nil
}

func (s *SubscriptionManager) CheckFeature(userID string, feature string) bool {
    tier := s.GetUserTier(userID)
    for _, f := range TierConfigs[tier].Features {
        if f == feature {
            return true
        }
    }
    return false
}
