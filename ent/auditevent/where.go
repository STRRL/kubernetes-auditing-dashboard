// Code generated by ent, DO NOT EDIT.

package auditevent

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldID, id))
}

// Level applies equality check predicate on the "level" field. It's identical to LevelEQ.
func Level(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldLevel, v))
}

// AuditID applies equality check predicate on the "auditID" field. It's identical to AuditIDEQ.
func AuditID(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldAuditID, v))
}

// Verb applies equality check predicate on the "verb" field. It's identical to VerbEQ.
func Verb(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldVerb, v))
}

// UserAgent applies equality check predicate on the "userAgent" field. It's identical to UserAgentEQ.
func UserAgent(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldUserAgent, v))
}

// Raw applies equality check predicate on the "raw" field. It's identical to RawEQ.
func Raw(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldRaw, v))
}

// RequestTimestamp applies equality check predicate on the "requestTimestamp" field. It's identical to RequestTimestampEQ.
func RequestTimestamp(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldRequestTimestamp, v))
}

// StageTimestamp applies equality check predicate on the "stageTimestamp" field. It's identical to StageTimestampEQ.
func StageTimestamp(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldStageTimestamp, v))
}

// LevelEQ applies the EQ predicate on the "level" field.
func LevelEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldLevel, v))
}

// LevelNEQ applies the NEQ predicate on the "level" field.
func LevelNEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldLevel, v))
}

// LevelIn applies the In predicate on the "level" field.
func LevelIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldLevel, vs...))
}

// LevelNotIn applies the NotIn predicate on the "level" field.
func LevelNotIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldLevel, vs...))
}

// LevelGT applies the GT predicate on the "level" field.
func LevelGT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldLevel, v))
}

// LevelGTE applies the GTE predicate on the "level" field.
func LevelGTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldLevel, v))
}

// LevelLT applies the LT predicate on the "level" field.
func LevelLT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldLevel, v))
}

// LevelLTE applies the LTE predicate on the "level" field.
func LevelLTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldLevel, v))
}

// LevelContains applies the Contains predicate on the "level" field.
func LevelContains(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContains(FieldLevel, v))
}

// LevelHasPrefix applies the HasPrefix predicate on the "level" field.
func LevelHasPrefix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasPrefix(FieldLevel, v))
}

// LevelHasSuffix applies the HasSuffix predicate on the "level" field.
func LevelHasSuffix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasSuffix(FieldLevel, v))
}

// LevelEqualFold applies the EqualFold predicate on the "level" field.
func LevelEqualFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEqualFold(FieldLevel, v))
}

// LevelContainsFold applies the ContainsFold predicate on the "level" field.
func LevelContainsFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContainsFold(FieldLevel, v))
}

// AuditIDEQ applies the EQ predicate on the "auditID" field.
func AuditIDEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldAuditID, v))
}

// AuditIDNEQ applies the NEQ predicate on the "auditID" field.
func AuditIDNEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldAuditID, v))
}

// AuditIDIn applies the In predicate on the "auditID" field.
func AuditIDIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldAuditID, vs...))
}

// AuditIDNotIn applies the NotIn predicate on the "auditID" field.
func AuditIDNotIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldAuditID, vs...))
}

// AuditIDGT applies the GT predicate on the "auditID" field.
func AuditIDGT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldAuditID, v))
}

// AuditIDGTE applies the GTE predicate on the "auditID" field.
func AuditIDGTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldAuditID, v))
}

// AuditIDLT applies the LT predicate on the "auditID" field.
func AuditIDLT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldAuditID, v))
}

// AuditIDLTE applies the LTE predicate on the "auditID" field.
func AuditIDLTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldAuditID, v))
}

// AuditIDContains applies the Contains predicate on the "auditID" field.
func AuditIDContains(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContains(FieldAuditID, v))
}

// AuditIDHasPrefix applies the HasPrefix predicate on the "auditID" field.
func AuditIDHasPrefix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasPrefix(FieldAuditID, v))
}

// AuditIDHasSuffix applies the HasSuffix predicate on the "auditID" field.
func AuditIDHasSuffix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasSuffix(FieldAuditID, v))
}

// AuditIDEqualFold applies the EqualFold predicate on the "auditID" field.
func AuditIDEqualFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEqualFold(FieldAuditID, v))
}

// AuditIDContainsFold applies the ContainsFold predicate on the "auditID" field.
func AuditIDContainsFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContainsFold(FieldAuditID, v))
}

// VerbEQ applies the EQ predicate on the "verb" field.
func VerbEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldVerb, v))
}

// VerbNEQ applies the NEQ predicate on the "verb" field.
func VerbNEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldVerb, v))
}

// VerbIn applies the In predicate on the "verb" field.
func VerbIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldVerb, vs...))
}

// VerbNotIn applies the NotIn predicate on the "verb" field.
func VerbNotIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldVerb, vs...))
}

// VerbGT applies the GT predicate on the "verb" field.
func VerbGT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldVerb, v))
}

// VerbGTE applies the GTE predicate on the "verb" field.
func VerbGTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldVerb, v))
}

// VerbLT applies the LT predicate on the "verb" field.
func VerbLT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldVerb, v))
}

// VerbLTE applies the LTE predicate on the "verb" field.
func VerbLTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldVerb, v))
}

// VerbContains applies the Contains predicate on the "verb" field.
func VerbContains(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContains(FieldVerb, v))
}

// VerbHasPrefix applies the HasPrefix predicate on the "verb" field.
func VerbHasPrefix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasPrefix(FieldVerb, v))
}

// VerbHasSuffix applies the HasSuffix predicate on the "verb" field.
func VerbHasSuffix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasSuffix(FieldVerb, v))
}

// VerbEqualFold applies the EqualFold predicate on the "verb" field.
func VerbEqualFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEqualFold(FieldVerb, v))
}

// VerbContainsFold applies the ContainsFold predicate on the "verb" field.
func VerbContainsFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContainsFold(FieldVerb, v))
}

// UserAgentEQ applies the EQ predicate on the "userAgent" field.
func UserAgentEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldUserAgent, v))
}

// UserAgentNEQ applies the NEQ predicate on the "userAgent" field.
func UserAgentNEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldUserAgent, v))
}

// UserAgentIn applies the In predicate on the "userAgent" field.
func UserAgentIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldUserAgent, vs...))
}

// UserAgentNotIn applies the NotIn predicate on the "userAgent" field.
func UserAgentNotIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldUserAgent, vs...))
}

// UserAgentGT applies the GT predicate on the "userAgent" field.
func UserAgentGT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldUserAgent, v))
}

// UserAgentGTE applies the GTE predicate on the "userAgent" field.
func UserAgentGTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldUserAgent, v))
}

// UserAgentLT applies the LT predicate on the "userAgent" field.
func UserAgentLT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldUserAgent, v))
}

// UserAgentLTE applies the LTE predicate on the "userAgent" field.
func UserAgentLTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldUserAgent, v))
}

// UserAgentContains applies the Contains predicate on the "userAgent" field.
func UserAgentContains(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContains(FieldUserAgent, v))
}

// UserAgentHasPrefix applies the HasPrefix predicate on the "userAgent" field.
func UserAgentHasPrefix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasPrefix(FieldUserAgent, v))
}

// UserAgentHasSuffix applies the HasSuffix predicate on the "userAgent" field.
func UserAgentHasSuffix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasSuffix(FieldUserAgent, v))
}

// UserAgentEqualFold applies the EqualFold predicate on the "userAgent" field.
func UserAgentEqualFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEqualFold(FieldUserAgent, v))
}

// UserAgentContainsFold applies the ContainsFold predicate on the "userAgent" field.
func UserAgentContainsFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContainsFold(FieldUserAgent, v))
}

// RawEQ applies the EQ predicate on the "raw" field.
func RawEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldRaw, v))
}

// RawNEQ applies the NEQ predicate on the "raw" field.
func RawNEQ(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldRaw, v))
}

// RawIn applies the In predicate on the "raw" field.
func RawIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldRaw, vs...))
}

// RawNotIn applies the NotIn predicate on the "raw" field.
func RawNotIn(vs ...string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldRaw, vs...))
}

// RawGT applies the GT predicate on the "raw" field.
func RawGT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldRaw, v))
}

// RawGTE applies the GTE predicate on the "raw" field.
func RawGTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldRaw, v))
}

// RawLT applies the LT predicate on the "raw" field.
func RawLT(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldRaw, v))
}

// RawLTE applies the LTE predicate on the "raw" field.
func RawLTE(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldRaw, v))
}

// RawContains applies the Contains predicate on the "raw" field.
func RawContains(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContains(FieldRaw, v))
}

// RawHasPrefix applies the HasPrefix predicate on the "raw" field.
func RawHasPrefix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasPrefix(FieldRaw, v))
}

// RawHasSuffix applies the HasSuffix predicate on the "raw" field.
func RawHasSuffix(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldHasSuffix(FieldRaw, v))
}

// RawEqualFold applies the EqualFold predicate on the "raw" field.
func RawEqualFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEqualFold(FieldRaw, v))
}

// RawContainsFold applies the ContainsFold predicate on the "raw" field.
func RawContainsFold(v string) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldContainsFold(FieldRaw, v))
}

// RequestTimestampEQ applies the EQ predicate on the "requestTimestamp" field.
func RequestTimestampEQ(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldRequestTimestamp, v))
}

// RequestTimestampNEQ applies the NEQ predicate on the "requestTimestamp" field.
func RequestTimestampNEQ(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldRequestTimestamp, v))
}

// RequestTimestampIn applies the In predicate on the "requestTimestamp" field.
func RequestTimestampIn(vs ...time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldRequestTimestamp, vs...))
}

// RequestTimestampNotIn applies the NotIn predicate on the "requestTimestamp" field.
func RequestTimestampNotIn(vs ...time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldRequestTimestamp, vs...))
}

// RequestTimestampGT applies the GT predicate on the "requestTimestamp" field.
func RequestTimestampGT(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldRequestTimestamp, v))
}

// RequestTimestampGTE applies the GTE predicate on the "requestTimestamp" field.
func RequestTimestampGTE(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldRequestTimestamp, v))
}

// RequestTimestampLT applies the LT predicate on the "requestTimestamp" field.
func RequestTimestampLT(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldRequestTimestamp, v))
}

// RequestTimestampLTE applies the LTE predicate on the "requestTimestamp" field.
func RequestTimestampLTE(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldRequestTimestamp, v))
}

// StageTimestampEQ applies the EQ predicate on the "stageTimestamp" field.
func StageTimestampEQ(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldEQ(FieldStageTimestamp, v))
}

// StageTimestampNEQ applies the NEQ predicate on the "stageTimestamp" field.
func StageTimestampNEQ(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNEQ(FieldStageTimestamp, v))
}

// StageTimestampIn applies the In predicate on the "stageTimestamp" field.
func StageTimestampIn(vs ...time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldIn(FieldStageTimestamp, vs...))
}

// StageTimestampNotIn applies the NotIn predicate on the "stageTimestamp" field.
func StageTimestampNotIn(vs ...time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldNotIn(FieldStageTimestamp, vs...))
}

// StageTimestampGT applies the GT predicate on the "stageTimestamp" field.
func StageTimestampGT(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGT(FieldStageTimestamp, v))
}

// StageTimestampGTE applies the GTE predicate on the "stageTimestamp" field.
func StageTimestampGTE(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldGTE(FieldStageTimestamp, v))
}

// StageTimestampLT applies the LT predicate on the "stageTimestamp" field.
func StageTimestampLT(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLT(FieldStageTimestamp, v))
}

// StageTimestampLTE applies the LTE predicate on the "stageTimestamp" field.
func StageTimestampLTE(v time.Time) predicate.AuditEvent {
	return predicate.AuditEvent(sql.FieldLTE(FieldStageTimestamp, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.AuditEvent) predicate.AuditEvent {
	return predicate.AuditEvent(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.AuditEvent) predicate.AuditEvent {
	return predicate.AuditEvent(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.AuditEvent) predicate.AuditEvent {
	return predicate.AuditEvent(func(s *sql.Selector) {
		p(s.Not())
	})
}
